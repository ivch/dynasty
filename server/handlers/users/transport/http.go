package transport

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/handlers/users/errs"
	"github.com/ivch/dynasty/server/middlewares"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Service interface {
	Register(ctx context.Context, req *users.User) (*users.User, error)
	// UserByPhoneAndPassword(ctx context.Context, phone, password string) (*entities.User, error)
	UserByID(ctx context.Context, id uint) (*users.User, error)
	AddFamilyMember(ctx context.Context, r *users.User) (*users.User, error)
	// ListFamilyMembers(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error)
	// DeleteFamilyMember(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error
}

type HTTPTransport struct {
	svc       Service
	log       logger.Logger
	router    chi.Router
	sanitizer *bluemonday.Policy
}

func (h *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewHTTPTransport returns a new instance of HTTPTransport.
func NewHTTPTransport(log logger.Logger, svc Service, p *bluemonday.Policy, mdl ...func(http.Handler) http.Handler) http.Handler {
	h := &HTTPTransport{log: log, router: chi.NewRouter().With(mdl...), svc: svc, sanitizer: p}
	h.attachRoutes()
	return h
}

func (h *HTTPTransport) attachRoutes() {
	h.router.Get("/v1/user", h.UserByID)
	h.router.Post("/v1/register", h.Register)
}

func (h *HTTPTransport) Register(w http.ResponseWriter, r *http.Request) {
	var req userRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}
	req.Sanitize(h.sanitizer)

	if err := validateRegisterRequest(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	user := users.User{
		Apartment:  req.Apartment,
		Email:      req.Email,
		Password:   req.Password,
		Phone:      req.Phone,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		BuildingID: req.BuildingID,
		EntryID:    req.EntryID,
		Active:     false,
		RegCode:    req.Code,
	}

	u, err := h.svc.Register(r.Context(), &user)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := userRegisterResponse{
		ID:    u.ID,
		Phone: u.Phone,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) UserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.UserByID(r.Context(), userID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	if res == nil {
		h.sendError(w, http.StatusNotFound, errs.UserNotFound)
		return
	}

	result := userByIDResponse{
		ID:        res.ID,
		Apartment: res.Apartment,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
		Building:  &res.Building,
		Entry:     &res.Entry,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func getUserID(ctx context.Context) (uint, error) {
	idStr, ok := middlewares.UserIDFromContext(ctx)
	if !ok {
		return 0, errs.EmptyUserID
	}

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, errs.BadUserID
	}

	if userID == 0 {
		return 0, errs.BadUserID
	}

	return uint(userID), nil
}

func (h *HTTPTransport) sendHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Debug("failed to send response error: %w", err)
	}
}

func (h *HTTPTransport) sendError(w http.ResponseWriter, httpCode int, error error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	res := errorResponse{
		ErrorCode: errs.Code(error),
		Error:     error.Error(),
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		h.log.Debug("failed to send response error: %w", err)
	}
}

func validateRegisterRequest(r *userRegisterRequest) error {
	if len(r.Password) < 6 {
		return errs.PasswordTooShort
	}

	if len(r.Phone) < 12 || len(r.Phone) > 13 {
		return errs.PhoneWrongLength
	}

	if _, err := strconv.ParseFloat(r.Phone, 64); err != nil {
		return errs.PhoneWrongChars
	}

	if len(r.FirstName) == 0 {
		return errs.FNameLength
	}

	if len(r.LastName) == 0 {
		return errs.LNameLength
	}

	if r.BuildingID == 0 {
		return errs.BuildingEmpty
	}

	if r.EntryID == 0 {
		return errs.EntryEmpty
	}

	if r.Apartment == 0 {
		return errs.ApartmentEmpty
	}

	if len(r.Email) == 0 {
		return errs.EmailEmpty
	}

	if !isEmailValid(r.Email) {
		return errs.EmailInvalid
	}

	if len(r.Code) < 5 {
		return errs.RegCodeLength
	}

	return nil
}

// isEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
