package requests

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
)

type loggingMiddleware struct {
	log  *zerolog.Logger
	next Service
}

func newLoggingMiddleware(log *zerolog.Logger, svc Service) Service {
	return loggingMiddleware{
		log:  log,
		next: svc,
	}
}

func (mw loggingMiddleware) Create(ctx context.Context, req *dto.RequestCreateRequest) (*dto.RequestCreateResponse, error) {
	res, err := mw.next.Create(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) My(ctx context.Context, req *dto.RequestMyRequest) (*dto.RequestMyResponse, error) {
	res, err := mw.next.My(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) Update(ctx context.Context, req *dto.RequestUpdateRequest) error {
	err := mw.next.Update(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return err
}

func (mw loggingMiddleware) Delete(ctx context.Context, req *dto.RequestByID) error {
	err := mw.next.Delete(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return err
}

func (mw loggingMiddleware) Get(ctx context.Context, req *dto.RequestByID) (*dto.RequestByIDResponse, error) {
	res, err := mw.next.Get(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) GuardRequestList(ctx context.Context, req *dto.GuardListRequest) (*dto.RequestMyResponse, error) {
	res, err := mw.next.GuardRequestList(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) GuardUpdateRequest(ctx context.Context, req *dto.GuardUpdateRequest) error {
	err := mw.next.GuardUpdateRequest(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return err
}
