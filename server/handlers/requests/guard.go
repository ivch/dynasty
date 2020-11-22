package requests

import (
	"context"
)

func (s *Service) GuardRequestList(_ context.Context, r *RequestListFilter) ([]*Request, int, error) {
	reqs, err := s.repo.ListForGuard(r)
	if err != nil {
		return nil, 0, err
	}

	cnt, err := s.repo.CountForGuard(r)
	if err != nil {
		return nil, 0, err
	}

	for i := range reqs {
		reqs[i].ImagesURL = make([]map[string]string, len(reqs[i].Images))
		for j := range reqs[i].Images {
			reqs[i].ImagesURL[j] = s.buildImageURL(reqs[i].Images[j])
		}
	}

	return reqs, cnt, nil
}

func (s *Service) GuardUpdateRequest(_ context.Context, r *Request) error {
	return s.repo.UpdateForGuard(r.ID, r.Status)
}
