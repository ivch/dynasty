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
		decodeUserByIDRequest,
		encodeHTTPResponse,
		options...))

	return r
}

func decodeUserByIDRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	idStr, ok := middleware.UserIDFromContext(ctx)

	if !ok || idStr == "" || idStr == "0" {
		return "", errors.New("empty id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return "", err
	}

	return uint(id), nil
}

func decodeRegisterRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var req userRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, errors.New("failed to decode request")
		}

		if err := validator.New().Struct(&req); err != nil {
			return nil, err
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
