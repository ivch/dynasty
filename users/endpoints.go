package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeRegisterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Register(ctx, request.(*userRegisterRequest))
	}
}

func makeUserByPhoneAndPasswordEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.UserByPhoneAndPassword(ctx, request.(*userByPhoneAndPasswordRequest))
	}
}

func makeUserByIDRequest(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.UserByID(ctx, request.(int))
	}
}
