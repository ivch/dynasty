package requests

import (
	"context"

	"github.com/rs/zerolog"
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

func (mw loggingMiddleware) Create(ctx context.Context, req *createRequest) (*createResponse, error) {
	res, err := mw.next.Create(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}
