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

func (s *Service) registerFamilyMember(_ context.Context, r *User, u *User) (*User, error) {
	if r.RegCode != u.RegCode {
		return nil, errs.RegCodeWrong
	}

	if u.Active {
		return nil, errs.FamilyMemberAlreadyRegistered
	}

	parent, err := s.repo.GetUserByID(*u.ParentID)
	if err != nil {
		s.log.Error("error getting parent user: %w", err)
		return nil, err
	}

	if parent.BuildingID != r.BuildingID || parent.EntryID != r.EntryID || parent.Apartment != r.Apartment {
		return nil, errs.FamilyMemberWrongAddress
	}

	pwd, err := hashAndSalt(r.Password)
	if err != nil {
		s.log.Error("error hashing password: %w", err)
		return nil, err
	}

	u.FirstName = r.FirstName
	u.LastName = r.LastName
	u.Email = r.Email
	u.Active = true
	u.Password = pwd

	if err := s.repo.UpdateUser(u); err != nil {
		return nil, err
	}

	return u, nil
}
