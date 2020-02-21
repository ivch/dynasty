package requests

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
	errEmptyID        = errors.New("empty id")
	errBadID          = errors.New("bad id")
	errEmptyUserID    = errors.New("empty user id")
	errBadUserID      = errors.New("bad user id")
	errBadRequest     = errors.New("failed to decode request")
	errInvalidRequest = errors.New("request validation error")
)

func New(log *zerolog.Logger, repo requestsRepository) (http.Handler, Service) {
	svc := newService(log, repo)
	h := newHTTPHandler(log, svc)
	return h, svc
}

func newHTTPHandler(log *zerolog.Logger, svc Service) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerBefore(middleware.UserIDToCTX),
	}

	r := chi.NewRouter()

	r.Method("POST", "/v1/request", httptransport.NewServer(
		makeCreateEndpoint(svc),
		decodeCreateRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("PUT", "/v1/request/{id}", httptransport.NewServer(
		makeUpdateEndpoint(svc),
		decodeUpdateRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/my", httptransport.NewServer(
		makeMyRequestsEndpoint(svc),
		decodeMyRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("DELETE", "/v1/request/{id}", httptransport.NewServer(
		makeDeleteEndpoint(svc),
		decodeByIDRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/request/{id}", httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeByIDRequest(log),
		encodeHTTPResponse,
		options...))

	// guard
	r.Method("GET", "/v1/guard/list", httptransport.NewServer(
		makeGuardRequestListEndpoint(svc),
		decodeGuardListRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("PUT", "/v1/guard/request/{id}", httptransport.NewServer(
		makeGuardUpdateRequest(svc),
		decodeGuardUpdateRequest(log),
		encodeHTTPResponse,
		options...))

	return r
}

func decodeGuardUpdateRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var req dto.GuardUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg(errBadRequest.Error())
			return nil, errBadRequest
		}

		idStr := chi.URLParam(r, "id")
		if idStr == "" || idStr == "0" {
			log.Error().Msg(errEmptyUserID.Error())
			return "", errEmptyID
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errBadID.Error())
			return "", errBadID
		}

		req.ID = uint(id)

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, errInvalidRequest
		}

		return &req, nil
	}
}

func decodeGuardListRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		offset, limit, err := parsePaginationRequest(r)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errBadRequest
		}

		req := dto.GuardListRequest{
			Offset: offset,
			Limit:  limit,
			Status: r.URL.Query().Get("status"),
		}

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, errInvalidRequest
		}

		return &req, nil
	}
}

func decodeByIDRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" || idStr == "0" {
			log.Error().Msg(errEmptyUserID.Error())
			return "", errEmptyID
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errBadID.Error())
			return "", errBadID
		}

		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, err
		}

		req := dto.RequestByID{
			UserID: userID,
			ID:     uint(id),
		}

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg(errInvalidRequest.Error())
			return nil, errInvalidRequest
		}

		return &req, nil
	}
}

func decodeUpdateRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var req dto.RequestUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg(errBadRequest.Error())
			return nil, errBadRequest
		}

		idStr := chi.URLParam(r, "id")
		if idStr == "" || idStr == "0" {
			log.Error().Msg(errEmptyUserID.Error())
			return nil, errEmptyID
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errBadID.Error())
			return nil, errBadID
		}

		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, err
		}

		req.ID = uint(id)
		req.UserID = userID

		return &req, nil
	}
}

func decodeMyRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, err
		}

		offset, limit, err := parsePaginationRequest(r)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errBadRequest
		}

		req := &dto.RequestMyRequest{
			UserID: userID,
			Offset: offset,
			Limit:  limit,
		}

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, errInvalidRequest
		}

		return req, nil
	}
}

func decodeCreateRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var req dto.RequestCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, errBadRequest
		}

		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, err
		}

		req.UserID = userID
		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, errInvalidRequest
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
	case errors.Is(err, errEmptyID):
		fallthrough
	case errors.Is(err, errBadID):
		fallthrough
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

func getUserID(ctx context.Context) (uint, error) {
	idStr, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return 0, errEmptyUserID
	}

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, errBadUserID
	}

	if userID == 0 {
		return 0, errBadUserID
	}

	return uint(userID), nil
}

func parsePaginationRequest(r *http.Request) (uint, uint, error) {
	_offset := r.URL.Query().Get("offset")
	_limit := r.URL.Query().Get("limit")

	if _offset == "" {
		return 0, 0, errors.New("empty offset")
	}

	if _limit == "" {
		return 0, 0, errors.New("empty limit")
	}

	offset, err := strconv.ParseUint(_offset, 10, 32)
	if err != nil {
		return 0, 0, errors.New("bad offset")
	}

	limit, err := strconv.ParseUint(_limit, 10, 32)
	if err != nil {
		return 0, 0, errors.New("bad limit")
	}

	if limit == 0 {
		return 0, 0, errors.New("limit should be grater then 0")
	}

	if limit > 200 {
		return 0, 0, errors.New("limit should less or equal 200")
	}

	return uint(offset), uint(limit), nil
}
