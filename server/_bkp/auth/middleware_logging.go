package auth

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

func (mw loggingMiddleware) Login(ctx context.Context, req *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
	res, err := mw.next.Login(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) Refresh(ctx context.Context, req *dto.AuthRefreshTokenRequest) (*dto.AuthLoginResponse, error) {
	res, err := mw.next.Refresh(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) Gwfa(token string) (uint, error) {
	id, err := mw.next.Gwfa(token)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return id, err
}
