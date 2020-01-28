package requests

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Create(ctx, request.(*createRequest))
	}
}
