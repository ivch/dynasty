package requests

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/ivch/dynasty/models/dto"
)

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Create(ctx, request.(*dto.RequestCreateRequest))
	}
}

func makeMyRequestsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.My(ctx, request.(*dto.RequestMyRequest))
	}
}

func makeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, svc.Update(ctx, request.(*dto.RequestUpdateRequest))
	}
}

func makeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, svc.Delete(ctx, request.(*dto.RequestByID))
	}
}

func makeGetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Get(ctx, request.(*dto.RequestByID))
	}
}
