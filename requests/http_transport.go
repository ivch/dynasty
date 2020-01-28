package requests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

const (
	userIDHeader = "X-Auth-User"
	userIDCtxKey = "userID"
)

func NewHTTPHandler(svc Service, log *zerolog.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerBefore(userIDToCTX),
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(accessLogMiddleware(log))

		r.Method("POST", "/requests/v1/request", httptransport.NewServer(
			makeCreateEndpoint(svc),
			decodeCreateRequest(log),
			encodeHTTPResponse,
			options...))
	})

	return r
}

func decodeCreateRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, err
		}

		userIDStr, ok := ctx.Value(userIDCtxKey).(string)
		if !ok {
			return nil, errors.New("user id is required")
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return nil, errors.New("wrong user id")
		}

		req.UserID = uint(userID)
		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, err
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

var userIDToCTX = func(ctx context.Context, r *http.Request) context.Context {
	if v := r.Header.Get(userIDHeader); v != "" {
		ctx = context.WithValue(ctx, userIDCtxKey, v)
	}

	return ctx
}

var accessLogMiddleware = func(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			duration := time.Since(start)
			log.Info().
				Str("tag", "http_log").
				Str("remote_addr", r.RemoteAddr).
				Str("request_method", r.Method).
				Str("request_uri", r.RequestURI).
				Int("response_code", ww.Status()).
				Dur("duration", duration).
				Msg("request")
		})
	}
}
