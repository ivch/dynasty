package users

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

func (mw loggingMiddleware) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	res, err := mw.next.Register(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) UserByPhoneAndPassword(ctx context.Context, phone, password string) (*dto.UserAuthResponse, error) {
	res, err := mw.next.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error) {
	res, err := mw.next.UserByID(ctx, id)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}
