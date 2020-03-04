package dictionaries

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

func (mw loggingMiddleware) BuildingsList(ctx context.Context) (*dto.BuildingsDictionaryResposnse, error) {
	res, err := mw.next.BuildingsList(ctx)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}
	return res, err
}

func (mw loggingMiddleware) EntriesList(ctx context.Context, id uint) (*dto.EntriesDictionaryResponse, error) {
	res, err := mw.next.EntriesList(ctx, id)
	if err != nil {
		mw.log.Error().Msg(err.Error())
	}
	return res, err
}
