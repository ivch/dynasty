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

func (mw loggingMiddleware) My(ctx context.Context, req *myRequest) (*myResponse, error) {
	res, err := mw.next.My(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}

func (mw loggingMiddleware) Update(ctx context.Context, req *updateRequest) error {
	err := mw.next.Update(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return err
}

func (mw loggingMiddleware) Delete(ctx context.Context, req *byIDRequest) error {
	err := mw.next.Delete(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return err
}

func (mw loggingMiddleware) Get(ctx context.Context, req *byIDRequest) (*getResponse, error) {
	res, err := mw.next.Get(ctx, req)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}

	return res, err
}
