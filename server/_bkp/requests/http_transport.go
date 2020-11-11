package requests

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/server/middlewares"
)

var (
	errEmptyID             = errors.New("empty id")
	errBadID               = errors.New("bad id")
	errEmptyUserID         = errors.New("empty user id")
	errBadUserID           = errors.New("bad user id")
	errBadRequest          = errors.New("failed to decode request")
	errInvalidRequest      = errors.New("request validation error")
	errInternalServerError = errors.New("request failed")
	errNoFile              = errors.New("error retrieving the file")
	errFileWrongType       = errors.New("wrong filetype")
	errFileIsTooBig        = errors.New("too big file")
	errTooMuchFiles        = errors.New("only allowed 3 images per request")
)

func New(log *zerolog.Logger, repo requestsRepository, p *bluemonday.Policy, s3Client s3Client, s3Space, cdnHost string) (http.Handler, Service) {
	svc := newService(log, repo, s3Client, s3Space, cdnHost)
	h := newHTTPHandler(log, svc, p)
	return h, svc
}

func newHTTPHandler(log *zerolog.Logger, svc Service, p *bluemonday.Policy) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerBefore(middlewares.UserIDToCTX),
	}

	r := chi.NewRouter()

	r.Method("POST", "/v1/request", httptransport.NewServer(
		makeCreateEndpoint(svc),
		decodeCreateRequest(log, p),
		encodeHTTPResponse,
		options...))

	r.Method("PUT", "/v1/request/{id}", httptransport.NewServer(
		makeUpdateEndpoint(svc),
		decodeUpdateRequest(log, p),
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

	r.Method("POST", "/v1/request/{id}/file", httptransport.NewServer(
		makeUploadImageEndpoint(svc),
		decodeUploadImageRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("DELETE", "/v1/request/{id}/file", httptransport.NewServer(
		makeDeleteImageEndpoint(svc),
		decodeDeleteImageRequest(log),
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

func decodeDeleteImageRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, errBadUserID
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

		req := dto.DeleteImageRequest{
			UserID:    userID,
			RequestID: uint(id),
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg(errBadRequest.Error())
			return nil, errBadRequest
		}

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, errInvalidRequest
		}

		return &req, nil
	}
}
func decodeUploadImageRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user id")
			return nil, errBadUserID
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

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Error().Err(err).Msg(err.Error())
			return nil, errInvalidRequest
		}

		file, header, err := r.FormFile("photo")
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return nil, errNoFile
		}

		if header.Size > (5 << 20) { // 5Mb
			log.Error().Err(err).Msg(errFileIsTooBig.Error())
			return nil, errFileIsTooBig
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error().Err(err).Msg("failed reading file content")
			return nil, errBadRequest
		}

		req := dto.UploadImageRequest{
			UserID:    userID,
			RequestID: uint(id),
			File:      fileBytes,
		}

		return &req, nil
	}
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

		req := dto.RequestListFilterRequest{
			Offset:    offset,
			Limit:     limit,
			Type:      r.URL.Query().Get("type"),
			Status:    r.URL.Query().Get("status"),
			Apartment: r.URL.Query().Get("apartment"),
			Place:     r.URL.Query().Get("place"),
		}

		if req.Type == "" {
			req.Type = "all" // nolint: goconst
		}

		if req.Place == "" {
			req.Place = "all" // nolint: goconst
		}

		if req.Status == "" {
			req.Status = "all" // nolint: goconst
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

		return &req, nil
	}
}

func decodeUpdateRequest(log *zerolog.Logger, p *bluemonday.Policy) httptransport.DecodeRequestFunc {
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
		req.Sanitize(p)
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

		req := dto.RequestListFilterRequest{
			UserID: userID,
			Offset: offset,
			Limit:  limit,
			Type:   "all",
			Status: "all",
		}

		return &req, nil
	}
}

func decodeCreateRequest(log *zerolog.Logger, p *bluemonday.Policy) httptransport.DecodeRequestFunc {
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
		req.Sanitize(p)
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
	case errors.Is(err, errNoFile):
		fallthrough
	case errors.Is(err, errFileIsTooBig):
		fallthrough
	case errors.Is(err, errFileWrongType):
		fallthrough
	case errors.Is(err, errInvalidRequest):
		status = http.StatusBadRequest
	default:
		err = errInternalServerError
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func getUserID(ctx context.Context) (uint, error) {
	idStr, ok := middlewares.UserIDFromContext(ctx)
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
