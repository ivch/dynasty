package users

import (
	"context"

	"github.com/ivch/dynasty/common/errs"
)

func (s *Service) AddFamilyMember(_ context.Context, r *User) (*User, error) {
	owner, err := s.repo.GetUserByID(*r.ParentID)
	if err != nil {
		return nil, err
	}

	members, err := s.repo.GetFamilyMembers(owner.ID)
	if err != nil {
		return nil, err
	}

	if len(members) >= s.membersLimit {
		return nil, errs.FamilyMembersLimitExceeded
	}

	u, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil {
		s.log.Error("error getting user by phone: %w", err)
		return nil, err
	}

	if u != nil {
		return nil, errs.FamilyMemberPhoneExists
	}

	regCode, err := s.repo.GetRegCode()
	if err != nil {
		s.log.Error("error getting reg code for user: %w", err)
		return nil, err
	}

	m := User{
		Phone:      r.Phone,
		Building:   owner.Building,
		Apartment:  owner.Apartment,
		Role:       defaultUserRole,
		BuildingID: owner.BuildingID,
		EntryID:    owner.EntryID,
		RegCode:    regCode,
		ParentID:   &owner.ID,
		Active:     false,
	}

	if err := s.repo.CreateUser(&m); err != nil {
		return nil, err
	}

	if err := s.repo.UseRegCode(regCode); err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *Service) ListFamilyMembers(_ context.Context, id uint) ([]*User, error) {
	list, err := s.repo.GetFamilyMembers(id)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) DeleteFamilyMember(_ context.Context, ownerID, memberID uint) error {
	member, err := s.repo.GetUserByID(memberID)
	if err != nil {
		return err
	}

	if *member.ParentID != ownerID {
		return errs.FamilyMemberWrongOwner
	}

	return s.repo.DeleteUser(member)
}

func (s *Service) registerFamilyMember(_ context.Context, request *User, member *User) (*User, error) {
	if request.RegCode != member.RegCode {
		return nil, errs.RegCodeWrong
	}

	if member.Active {
		return nil, errs.FamilyMemberAlreadyRegistered
	}

	parent, err := s.repo.GetUserByID(*member.ParentID)
	if err != nil {
		s.log.Error("error getting parent user: %w", err)
		return nil, err
	}

	if parent.BuildingID != request.BuildingID || parent.EntryID != request.EntryID || parent.Apartment != request.Apartment {
		return nil, errs.FamilyMemberWrongAddress
	}

	pwd, err := hashAndSalt(request.Password)
	if err != nil {
		s.log.Error("error hashing password: %w", err)
		return nil, err
	}

	member.Active = true

	update := UserUpdate{
		ID:        member.ID,
		Email:     &request.Email,
		Password:  &pwd,
		FirstName: &request.FirstName,
		LastName:  &request.LastName,
		Active:    &member.Active,
	}

	if err := s.repo.UpdateUser(&update); err != nil {
		return nil, err
	}

	return member, nil
}
