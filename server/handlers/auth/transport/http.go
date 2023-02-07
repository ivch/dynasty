package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/auth"
	"github.com/ivch/dynasty/server/middlewares"
)

var uuidRegexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

type AuthService interface {
	Login(ctx context.Context, phone, password string) (*auth.Tokens, error)
	Logout(ctx context.Context, id uint) error
	Refresh(ctx context.Context, token string) (*auth.Tokens, error)
	Gwfa(token string) (uint, error)
}

type HTTPTransport struct {
	svc    AuthService
	log    logger.Logger
	router chi.Router
}

func (h *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewHTTPTransport returns a new instance of HTTPTransport.
func NewHTTPTransport(log logger.Logger, svc AuthService, mdl ...func(http.Handler) http.Handler) http.Handler {
	h := &HTTPTransport{log: log, router: chi.NewRouter().With(mdl...), svc: svc}
	h.attachRoutes()
	return h
}

func (h *HTTPTransport) attachRoutes() {
	h.router.Post("/v1/login", h.Login)
	h.router.Get("/v1/logout", h.Logout)
	h.router.Post("/v1/refresh", h.Refresh)
	h.router.Get("/v1/gwfa", h.Gwfa)
}

func (h *HTTPTransport) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	if err := validateLoginRequest(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	// ip, _, err := net.SplitHostPort(r.RemoteAddr)
	// if err != nil {
	//	log.Error().Err(err).Msg("failed to parse host")
	//	return nil, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	// }
	//
	// userIP := net.ParseIP(ip)
	// if userIP == nil {
	//	log.Error().Err(err).Msg("failed to parse ip")
	//	return nil, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	// }
	//
	// req.IP = userIP
	// req.Ua = r.Header.Get("User-Agent")

	res, err := h.svc.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := loginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) Refresh(w http.ResponseWriter, r *http.Request) {
	var req authRefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	if !uuidRegexp.MatchString(req.Token) {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	res, err := h.svc.Refresh(r.Context(), req.Token)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := loginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) Gwfa(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		h.sendError(w, http.StatusUnauthorized, errs.NoAuthHeader)
		return
	}

	id, err := h.svc.Gwfa(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	w.Header().Set("X-Auth-User", strconv.FormatUint(uint64(id), 10))

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	if err := h.svc.Logout(r.Context(), userID); err != nil {
		h.sendError(w, http.StatusInternalServerError, nil)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func validateLoginRequest(r *loginRequest) error {
	if len(r.Password) < 6 {
		return errs.PasswordTooShort
	}

	if len(r.Phone) < 12 || len(r.Phone) > 13 {
		return errs.PhoneWrongLength
	}

	if _, err := strconv.ParseFloat(r.Phone, 64); err != nil {
		return errs.PhoneWrongChars
	}

	return nil
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
