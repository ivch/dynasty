package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/ivch/dynasty/models/dto"
)

func makeRegisterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Register(ctx, request.(*dto.UserRegisterRequest))
	}
}

func makeUserByIDRequest(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.UserByID(ctx, request.(uint))
	}
}
