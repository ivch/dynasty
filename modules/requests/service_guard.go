package requests

import (
	"context"
	"fmt"

	"github.com/ivch/dynasty/models/dto"
)

func (s *service) GuardRequestList(_ context.Context, r *dto.RequestListFilterRequest) (*dto.RequestGuardListResponse, error) {
	reqs, err := s.repo.ListForGuard(r)
	if err != nil {
		return nil, err
	}

	cnt, err := s.repo.CountForGuard(r)
	if err != nil {
		return nil, err
	}

	data := make([]*dto.RequestForGuard, len(reqs))
	for i := range reqs {
		data[i] = &dto.RequestForGuard{
			ID:          reqs[i].ID,
			UserID:      reqs[i].UserID,
			Type:        reqs[i].Type,
			Time:        reqs[i].Time,
			Description: reqs[i].Description,
			Status:      reqs[i].Status,
			UserName:    fmt.Sprintf("%s %s", reqs[i].User.FirstName, reqs[i].User.LastName),
			Phone:       reqs[i].User.Phone,
			Address:     reqs[i].User.Building.Name + ", " + reqs[i].User.Entry.Name,
			Apartment:   reqs[i].User.Apartment,
			Images:      make([]map[string]string, len(reqs[i].Images)),
		}

		for j := range reqs[i].Images {
			data[i].Images[j] = s.buildImageURL(reqs[i].Images[j])
		}
	}

	return &dto.RequestGuardListResponse{
		Data:  data,
		Count: cnt,
	}, nil
}

func (s *service) GuardUpdateRequest(_ context.Context, r *dto.GuardUpdateRequest) error {
	return s.repo.UpdateForGuard(r.ID, r.Status)
}
