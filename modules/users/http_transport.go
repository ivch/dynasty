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
	errEmptyUserID                   = errors.New("empty user id")
	errBadUserID                     = errors.New("bad user id")
	errBadRequest                    = errors.New("failed to decode request")
	errInvalidRequest                = errors.New("request validation error")
	errMasterAccountExists           = errors.New("master account for this apt already exists")
	errFamilyMembersLimitExceeded    = errors.New("family members limit exceeded")
	errProvidedWrongRegCode          = errors.New("provided wrong reg code")
	errFamilyMemberAlreadyRegistered = errors.New("family member already registered")
)

func New(repo userRepository, verifyRegCode bool, membersLimit int, log *zerolog.Logger) (http.Handler, Service) {
	svc := newService(log, repo, verifyRegCode, membersLimit)
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
		makeUserByIDEndpoint(svc),
		decodeUserByIDRequest,
		encodeHTTPResponse,
		options...))

	r.Method("POST", "/v1/member", httptransport.NewServer(
		makeAddFamilyMemberEndpoint(svc),
		decodeAddFamilyMemberRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/members", httptransport.NewServer(
		makeListFamilyMembersEndpoint(svc),
		decodeListFamilyMembersRequest,
		encodeHTTPResponse,
		options...))

	r.Method("DELETE", "/v1/member/{id}", httptransport.NewServer(
		makeDeleteFamilyMemberEndpoint(svc),
		decodeDeleteFamilyMemberRequest(log),
		encodeHTTPResponse,
		options...))

	return r
}

func decodeDeleteFamilyMemberRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		ownerID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		memberID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg(errBadUserID.Error())
			return nil, errBadUserID
		}

		if memberID == 0 {
			log.Error().Msg(errEmptyUserID.Error())
			return nil, errEmptyUserID
		}

		if uint(memberID) == ownerID {
			return nil, errBadRequest
		}

		return &dto.DeleteFamilyMemberRequest{
			OwnerID:  ownerID,
			MemberID: uint(memberID),
		}, nil
	}
}

func decodeListFamilyMembersRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return getUserID(ctx)
}

func decodeAddFamilyMemberRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		id, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		var req dto.AddFamilyMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, errBadRequest
		}

		req.OwnerID = id

		if err := validator.New().Struct(&req); err != nil {
			return nil, errInvalidRequest
		}

		return &req, nil
	}
}

func decodeUserByIDRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return getUserID(ctx)
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
	var (
		status  = http.StatusInternalServerError
		message = err.Error()
	)

	switch {
	case errors.Is(err, errEmptyUserID):
		fallthrough
	case errors.Is(err, errBadUserID):
		fallthrough
	case errors.Is(err, errBadRequest):
		fallthrough
	case errors.Is(err, errInvalidRequest):
		fallthrough
	case errors.Is(err, errProvidedWrongRegCode):
		status = http.StatusBadRequest
	case errors.Is(err, errMasterAccountExists):
		status = http.StatusConflict
	default:
		message = "something went wrong"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
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
