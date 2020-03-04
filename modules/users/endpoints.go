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

func makeUserByIDEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.UserByID(ctx, request.(uint))
	}
}

func makeAddFamilyMemberEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.AddFamilyMember(ctx, request.(*dto.AddFamilyMemberRequest))
	}
}

func makeListFamilyMembersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return svc.ListFamilyMembers(ctx, request.(uint))
	}
}

func makeDeleteFamilyMemberEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, svc.DeleteFamilyMember(ctx, request.(*dto.DeleteFamilyMemberRequest))
	}
}
