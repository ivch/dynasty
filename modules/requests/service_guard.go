package requests

import (
	"context"
	"fmt"

	"github.com/ivch/dynasty/models/dto"
)

func (s *service) GuardRequestList(_ context.Context, r *dto.RequestListFilterRequest) ([]*dto.RequestForGuard, error) {
	reqs, err := s.repo.ListForGuard(r)
	if err != nil {
		return nil, err
	}

	res := make([]*dto.RequestForGuard, len(reqs))
	for i, v := range reqs {
		res[i] = &dto.RequestForGuard{
			ID:          v.ID,
			UserID:      v.UserID,
			Type:        v.Type,
			Time:        v.Time,
			Description: v.Description,
			Status:      v.Status,
			UserName:    fmt.Sprintf("%s %s", v.User.FirstName, v.User.LastName),
			Phone:       v.User.Phone,
			Address:     v.User.Building.Address,
			Apartment:   v.User.Apartment,
		}
	}
	return res, nil
}

func (s *service) GuardUpdateRequest(_ context.Context, r *dto.GuardUpdateRequest) error {
	return s.repo.UpdateForGuard(r.ID, r.Status)
}
