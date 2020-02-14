package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"

	"github.com/ivch/dynasty/common/middleware"
	"github.com/ivch/dynasty/models/dto"
)

var (
	errEmptyUserID    = errors.New("empty user id")
	errBadUserID      = errors.New("bad user id")
	errBadRequest     = errors.New("failed to decode request")
	errInvalidRequest = errors.New("request validation error")
)

func New(repo userRepository, verifyRegCode bool, log *zerolog.Logger) (http.Handler, Service) {
	svc := newService(log, repo, verifyRegCode)
	h := newHTTPHandler(log, svc)

	return h, svc
}

func newHTTPHandler(log *zerolog.Logger, svc Service) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerBefore(middleware.UserIDToCTX),
	}

	r := chi.NewRouter()
	r.Method("POST", "/v1/register", httptransport.NewServer(
		makeRegisterEndpoint(svc),
		decodeRegisterRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/user", httptransport.NewServer(
		makeUserByIDRequest(svc),
		decodeUserByIDRequest(log),
		encodeHTTPResponse,
		options...))

	return r
}

func decodeUserByIDRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, _ *http.Request) (interface{}, error) {
		idStr, ok := middleware.UserIDFromContext(ctx)

		if !ok || idStr == "" || idStr == "0" {
			return "", errEmptyUserID
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errBadUserID.Error())
			return "", errBadUserID
		}

		return uint(id), nil
	}
}

func decodeRegisterRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var req dto.UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, errBadRequest
		}

		if err := validator.New().Struct(&req); err != nil {
			return nil, errInvalidRequest
		}

		if req.Phone[0] == '+' {
			req.Phone = req.Phone[1:]
		}

		return &req, nil
	}
}

func encodeHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	var status int
	switch {
	case errors.Is(err, errEmptyUserID):
		fallthrough
	case errors.Is(err, errBadUserID):
		fallthrough
	case errors.Is(err, errBadRequest):
		fallthrough
	case errors.Is(err, errInvalidRequest):
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
