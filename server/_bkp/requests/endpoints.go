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
		return svc.My(ctx, request.(*dto.RequestListFilterRequest))
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

func makeGuardRequestListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.GuardRequestList(ctx, request.(*dto.RequestListFilterRequest))
	}
}
func makeGuardUpdateRequest(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, svc.GuardUpdateRequest(ctx, request.(*dto.GuardUpdateRequest))
	}
}

func makeUploadImageEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.UploadImage(ctx, request.(*dto.UploadImageRequest))
	}
}

func makeDeleteImageEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, svc.DeleteImage(ctx, request.(*dto.DeleteImageRequest))
	}
}
