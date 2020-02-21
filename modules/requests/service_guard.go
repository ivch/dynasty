package requests

import (
	"context"

	"github.com/ivch/dynasty/models/dto"
)

func (s *service) GuardRequestList(_ context.Context, r *dto.GuardListRequest) (*dto.RequestMyResponse, error) {
	reqs, err := s.repo.ListForGuard(r.Limit, r.Offset, r.Status)
	if err != nil {
		return nil, err
	}
	return &dto.RequestMyResponse{Data: reqs}, nil
}

func (s *service) GuardUpdateRequest(_ context.Context, r *dto.GuardUpdateRequest) error {
	return s.repo.UpdateForGuard(r.ID, r.Status)
}
