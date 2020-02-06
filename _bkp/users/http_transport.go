package users

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

		r.Method("POST", "/users/v1/register", httptransport.NewServer(
			makeRegisterEndpoint(svc),
			decodeRegisterRequest(log),
			encodeHTTPResponse,
			options...))

		r.Method("GET", "/users/v1/internal/user-phone-and-pass", httptransport.NewServer(
			makeUserByPhoneAndPasswordEndpoint(svc),
			decodeUserByPhoneAndPasswordRequest(log),
			encodeHTTPResponse,
			options...))

		r.Method("GET", "/users/v1/internal/user/{id}", httptransport.NewServer(
			makeUserByIDRequest(svc),
			decodeUserByIDRequest(true),
			encodeHTTPResponse,
			options...))

		r.Method("GET", "/users/v1/user", httptransport.NewServer(
			makeUserByIDRequest(svc),
			decodeUserByIDRequest(false),
			encodeHTTPResponse,
			options...))
	})

	return r
}

func decodeUserByIDRequest(internal bool) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		idStr := chi.URLParam(r, "id")
		if !internal {
			idStr = ctx.Value(userIDCtxKey).(string)
		}

		if idStr == "" || idStr == "0" {
			return "", errors.New("empty id")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func decodeUserByPhoneAndPasswordRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (request interface{}, err error) {
		req := userByPhoneAndPasswordRequest{
			Phone:    r.URL.Query().Get("phone"),
			Password: r.URL.Query().Get("pwd"),
		}
		if err := validator.New().Struct(&req); err != nil {
			log.Error().Msg(err.Error())
			return nil, err
		}

		return &req, nil
	}
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
