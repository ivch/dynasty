package users

import (
	"context"

	"github.com/ivch/dynasty/server/handlers/users/errs"
)

// func (s *service) AddFamilyMember(_ context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error) {
// 	owner, err := s.repo.GetUserByID(r.OwnerID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	members, err := s.repo.GetFamilyMembers(owner.ID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(members) >= s.membersLimit {
// 		return nil, errFamilyMembersLimitExceeded
// 	}
//
// 	u, err := s.repo.GetUserByPhone(r.Phone)
// 	if err != nil && err != entities.ErrUserNotFound {
// 		return nil, err
// 	}
//
// 	if u != nil {
// 		return nil, entities.ErrUserPhoneExists
// 	}
//
// 	regCode, err := s.repo.GetRegCode()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	m := entities.User{
// 		Phone:      r.Phone,
// 		Building:   owner.Building,
// 		Apartment:  owner.Apartment,
// 		Role:       entities.DefaultUserRole,
// 		BuildingID: owner.BuildingID,
// 		RegCode:    regCode,
// 		ParentID:   &owner.ID,
// 		Active:     false,
// 	}
//
// 	if err := s.repo.CreateUser(&m); err != nil {
// 		return nil, err
// 	}
//
// 	if err := s.repo.UseRegCode(regCode); err != nil {
// 		return nil, err
// 	}
//
// 	return &dto.AddFamilyMemberResponse{Code: regCode}, nil
// }
//
// func (s *service) ListFamilyMembers(_ context.Context, id uint) (*dto.ListFamilyMembersResponse, error) {
// 	list, err := s.repo.GetFamilyMembers(id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(list) == 0 {
// 		return nil, nil
// 	}
//
// 	members := make([]*dto.FamilyMember, len(list))
// 	for i, m := range list {
// 		members[i] = &dto.FamilyMember{
// 			ID:     m.ID,
// 			Name:   m.FirstName + " " + m.LastName,
// 			Phone:  m.Phone,
// 			Code:   m.RegCode,
// 			Active: m.Active,
// 		}
// 	}
//
// 	return &dto.ListFamilyMembersResponse{Data: members}, nil
// }
//
// func (s *service) DeleteFamilyMember(_ context.Context, r *dto.DeleteFamilyMemberRequest) error {
// 	member, err := s.repo.GetUserByID(r.MemberID)
// 	if err != nil {
// 		return err
// 	}
//
// 	if *member.ParentID != r.OwnerID {
// 		return errors.New("wrong owner or member id")
// 	}
//
// 	return s.repo.DeleteUser(member)
// }

func (s *service) registerFamilyMember(_ context.Context, r *User, u *User) (*User, error) {
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
