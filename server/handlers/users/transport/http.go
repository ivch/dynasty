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

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/middlewares"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type UsersService interface {
	Register(ctx context.Context, req *users.User) (*users.User, error)
	Update(ctx context.Context, req *users.UserUpdate) error
	UserByID(ctx context.Context, id uint) (*users.User, error)
	AddFamilyMember(ctx context.Context, r *users.User) (*users.User, error)
	ListFamilyMembers(ctx context.Context, id uint) ([]*users.User, error)
	DeleteFamilyMember(ctx context.Context, ownerID, memberID uint) error
	RecoveryCode(ctx context.Context, r *users.User) error
	ResetPassword(ctx context.Context, code string, r *users.UserUpdate) error
}

type HTTPTransport struct {
	svc       UsersService
	log       logger.Logger
	router    chi.Router
	sanitizer *bluemonday.Policy
}

func (h *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewHTTPTransport returns a new instance of HTTPTransport.
func NewHTTPTransport(log logger.Logger, svc UsersService, p *bluemonday.Policy, mdl ...func(http.Handler) http.Handler) http.Handler {
	h := &HTTPTransport{log: log, router: chi.NewRouter().With(mdl...), svc: svc, sanitizer: p}
	h.attachRoutes()
	return h
}

func (h *HTTPTransport) attachRoutes() {
	h.router.Get("/v1/user", h.UserByID)
	h.router.Put("/v1/user", h.Update)
	h.router.Post("/v1/register", h.Register)
	h.router.Post("/v1/member", h.AddFamilyMember)
	h.router.Get("/v1/members", h.FamilyMembersList)
	h.router.Delete("/v1/member/{id}", h.DeleteFamilyMember)
	h.router.Post("/v1/password-recovery", h.PasswordRecovery)
	h.router.Post("/v1/password-reset", h.PasswordReset)
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

func (h *HTTPTransport) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	var req userUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	req.Sanitize(h.sanitizer)

	data := users.UserUpdate{
		ID: userID,
	}

	if req.Email != nil {
		if !isEmailValid(*req.Email) {
			h.sendError(w, http.StatusBadRequest, errs.EmailInvalid)
			return
		}
		data.Email = req.Email
	}

	if req.FirstName != nil {
		data.FirstName = req.FirstName
	}

	if req.LastName != nil {
		data.LastName = req.LastName
	}

	if req.NewPassword != nil {
		if req.Password == nil {
			h.sendError(w, http.StatusBadRequest, errs.InvalidCredentials)
			return
		}

		if req.NewPasswordConfirm == nil {
			h.sendError(w, http.StatusBadRequest, errs.PasswordConfirmMismatch)
			return
		}

		if *req.NewPasswordConfirm != *req.NewPassword {
			h.sendError(w, http.StatusBadRequest, errs.PasswordConfirmMismatch)
			return
		}

		if err := validatePassword(*req.NewPassword); err != nil {
			h.sendError(w, http.StatusBadRequest, err)
			return
		}

		data.Password = req.Password
		data.NewPassword = req.NewPassword
	}

	if err := h.svc.Update(r.Context(), &data); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
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

	result := UserByIDResponse{
		ID:        res.ID,
		Apartment: res.Apartment,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
		Building:  &res.Building,
		Entry:     &res.Entry,
		Role:      res.Role,
		Active:    res.Active,
		ParentID:  res.ParentID,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) DeleteFamilyMember(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	memberID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.log.Error("wrong member id: %w", err)
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	if memberID == 0 {
		h.sendError(w, http.StatusBadRequest, errs.FamilyMemberBadID)
		return
	}

	if uint(memberID) == userID {
		h.sendError(w, http.StatusBadRequest, errs.FamilyMemberParentMismatch)
		return
	}

	if err := h.svc.DeleteFamilyMember(r.Context(), userID, uint(memberID)); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) FamilyMembersList(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.ListFamilyMembers(r.Context(), userID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	var result listFamilyMembersResponse
	result.Data = make([]*familyMember, len(res))
	if len(res) > 0 {
		for i, m := range res {
			result.Data[i] = &familyMember{
				ID:     m.ID,
				Name:   m.FirstName + " " + m.LastName,
				Phone:  m.Phone,
				Code:   m.RegCode,
				Active: m.Active,
			}
		}
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) AddFamilyMember(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	var req addFamilyMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	req.OwnerID = userID

	if err := validatePhone(req.Phone); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.PhoneWrongLength)
		return
	}

	u := users.User{
		Phone:    req.Phone,
		ParentID: &req.OwnerID,
	}

	res, err := h.svc.AddFamilyMember(r.Context(), &u)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := addFamilyMemberResponse{
		Code: res.RegCode,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) PasswordRecovery(w http.ResponseWriter, r *http.Request) {
	var req passwordRecoveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	if err := validatePhone(req.Phone); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.PhoneWrongLength)
		return
	}

	if !isEmailValid(req.Email) {
		h.sendError(w, http.StatusBadRequest, errs.EmailInvalid)
		return
	}

	req.Email = strings.ToLower(req.Email)

	if err := h.svc.RecoveryCode(r.Context(), &users.User{
		Email: req.Email,
		Phone: req.Phone,
	}); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	if len(req.Code) < 10 {
		h.sendError(w, http.StatusBadRequest, errs.BadRecoveryCode)
		return
	}

	var data users.UserUpdate
	if req.NewPassword == "" {
		h.sendError(w, http.StatusBadRequest, errs.EmptyPassword)
		return
	}

	if req.NewPasswordConfirm == "" {
		h.sendError(w, http.StatusBadRequest, errs.PasswordConfirmMismatch)
		return
	}

	if req.NewPasswordConfirm != req.NewPassword {
		h.sendError(w, http.StatusBadRequest, errs.PasswordConfirmMismatch)
		return
	}

	if err := validatePassword(req.NewPassword); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	data.NewPassword = &req.NewPassword

	if err := h.svc.ResetPassword(r.Context(), req.Code, &data); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
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

	if error == nil {
		error = errs.Generic
	}

	var (
		ru string
		ua string
	)

	if e, ok := error.(errs.SvcError); ok {
		ru, ua = e.Ru, e.Ua
	}

	res := errorResponse{
		ErrorCode: errs.Code(error),
		Error:     error.Error(),
		Ru:        ru,
		Ua:        ua,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		h.log.Debug("failed to send response error: %w", err)
	}
}

func validatePassword(p string) error {
	if len(p) < 6 {
		return errs.PasswordTooShort
	}
	return nil
}

func validatePhone(p string) error {
	if len(p) < 12 || len(p) > 13 {
		return errs.PhoneWrongLength
	}

	if _, err := strconv.ParseFloat(p, 64); err != nil {
		return errs.PhoneWrongChars
	}

	return nil
}

func validateRegisterRequest(r *userRegisterRequest) error {
	if err := validatePassword(r.Password); err != nil {
		return err
	}

	if err := validatePhone(r.Phone); err != nil {
		return err
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

	if r.Apartment > 1050 {
		return errs.AptNumberIsTooBig
	}

	if len(r.Email) == 0 {
		return errs.EmailEmpty
	}

	if !isEmailValid(r.Email) {
		return errs.EmailInvalid
	}

	return nil
}

// isEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
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
