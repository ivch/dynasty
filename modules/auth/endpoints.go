package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/ivch/dynasty/models/dto"
)

func makeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Login(ctx, request.(*dto.AuthLoginRequest))
	}
}

func makeRefreshEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.Refresh(ctx, request.(*dto.AuthRefreshTokenRequest))
	}
}
