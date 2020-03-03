package dictionaries

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/common/middleware"
)

var (
	errEmptyID = errors.New("empty id")
	errBadID   = errors.New("bad id")
)

func New(repo dictRepository, log *zerolog.Logger) (http.Handler, Service) {
	svc := newService(repo, log)
	h := newHTTPHandler(log, svc)

	return h, svc
}

func newHTTPHandler(logger *zerolog.Logger, svc Service) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerBefore(middleware.UserIDToCTX),
	}

	r := chi.NewRouter()
	r.Method("GET", "/v1/buildings", httptransport.NewServer(
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			return svc.BuildingsList(ctx)
		},
		func(_ context.Context, r *http.Request) (interface{}, error) { return nil, nil },
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/building/{id}/entries", httptransport.NewServer(
		func(ctx context.Context, r interface{}) (interface{}, error) {
			return svc.EntriesList(ctx, r.(uint))
		},
		decodeEntriesListRequest(logger),
		encodeHTTPResponse,
		options...))

	return r
}

func decodeEntriesListRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errEmptyID.Error())
			return nil, errEmptyID
		}

		if id == 0 {
			log.Error().Err(err).Msg(errBadID.Error())
			return nil, errBadID
		}

		return uint(id), nil
	}
}

func encodeHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	var (
		status  = http.StatusInternalServerError
		message = err.Error()
	)

	switch {
	case errors.Is(err, errBadID):
		fallthrough
	case errors.Is(err, errEmptyID):
		status = http.StatusBadRequest
	default:
		message = "something went wrong"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
