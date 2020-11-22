package transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/requests"
	"github.com/ivch/dynasty/server/middlewares"
)

type RequestsService interface {
	Create(ctx context.Context, r *requests.Request) (*requests.Request, error)
	Get(ctx context.Context, r *requests.Request) (*requests.Request, error)
	Update(ctx context.Context, r *requests.Request) error
	Delete(ctx context.Context, r *requests.Request) error
	My(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, error)

	UploadImage(ctx context.Context, r *requests.Image) (*requests.Image, error)
	DeleteImage(ctx context.Context, r *requests.Image) error

	GuardRequestList(ctx context.Context, r *requests.RequestListFilter) ([]*requests.Request, int, error)
	GuardUpdateRequest(ctx context.Context, r *requests.Request) error
}

type HTTPTransport struct {
	svc       RequestsService
	log       logger.Logger
	router    chi.Router
	sanitizer *bluemonday.Policy
}

func (h *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewHTTPTransport returns a new instance of HTTPTransport.
func NewHTTPTransport(log logger.Logger, svc RequestsService, p *bluemonday.Policy, mdl ...func(http.Handler) http.Handler) http.Handler {
	h := &HTTPTransport{log: log, router: chi.NewRouter().With(mdl...), svc: svc, sanitizer: p}
	h.attachRoutes()
	return h
}

func (h *HTTPTransport) attachRoutes() {
	h.router.Post("/v1/request", h.Create)
	h.router.Put("/v1/request/{id}", h.Update)
	h.router.Get("/v1/request/{id}", h.GetRequestByID)
	h.router.Delete("/v1/request/{id}", h.Delete)
	h.router.Get("/v1/my", h.ListByUser)

	h.router.Post("/v1/request/{id}/file", h.UploadFile)
	h.router.Delete("/v1/request/{id}/file", h.DeleteFile)

	h.router.Get("/v1/guard/list", h.GuardList)
	h.router.Put("/v1/guard/request/{id}", h.GuardUpdateRequest)
}

func (h *HTTPTransport) Create(w http.ResponseWriter, r *http.Request) {
	var req RequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	req.UserID = userID

	req.Sanitize(h.sanitizer)

	if err := validateCreateRequest(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	data := requests.Request{
		Type:        req.Type,
		UserID:      req.UserID,
		Time:        req.Time,
		Description: req.Description,
	}

	res, err := h.svc.Create(r.Context(), &data)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, RequestCreateResponse{ID: res.ID})
}

func (h *HTTPTransport) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	var req RequestUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	req.Sanitize(h.sanitizer)

	data := requests.Request{
		ID:     id,
		UserID: userID,
	}

	if req.Type != nil {
		data.Type = *req.Type
	}

	if req.Description != nil {
		data.Description = *req.Description
	}

	if req.Status != nil {
		data.Status = *req.Status
	}

	if req.Time != nil {
		data.Time = *req.Time
	}

	if err := h.svc.Update(r.Context(), &data); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	req := requests.Request{
		ID:     id,
		UserID: userID,
	}

	res, err := h.svc.Get(r.Context(), &req)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := RequestByIDResponse{
		ID:          res.ID,
		Type:        res.Type,
		UserID:      res.UserID,
		Time:        res.Time,
		Description: res.Description,
		Status:      res.Status,
		Images:      res.ImagesURL,
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	req := requests.Request{
		ID:     id,
		UserID: userID,
	}

	if err := h.svc.Delete(r.Context(), &req); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	offset, limit, err := parsePaginationRequest(r)
	if err != nil {
		h.log.Error("bad pagination request: %w", err)
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	req := requests.RequestListFilter{
		Type:   "all",
		Offset: offset,
		Limit:  limit,
		UserID: userID,
		Status: "all",
	}

	res, err := h.svc.My(r.Context(), &req)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]*RequestByIDResponse, len(res))
	for i := range res {
		result[i] = &RequestByIDResponse{
			ID:          res[i].ID,
			Type:        res[i].Type,
			UserID:      res[i].UserID,
			Time:        res[i].Time,
			Description: res[i].Description,
			Status:      res[i].Status,
			Images:      res[i].ImagesURL,
		}
	}

	h.sendHTTPResponse(r.Context(), w, ListByUserResponse{Data: result})
}

func (h *HTTPTransport) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error("error parsing file: %w", err)
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		h.log.Error("error reading file: %w", err)
		h.sendError(w, http.StatusBadRequest, errs.NoFile)
		return
	}

	if header.Size > (5 << 20) { // 5Mb
		h.sendError(w, http.StatusBadRequest, errs.FileIsTooBig)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		h.log.Error("failed reading file content: %w", err)
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	upload := requests.Image{
		UserID:    userID,
		RequestID: id,
		File:      fileBytes,
	}

	img, err := h.svc.UploadImage(r.Context(), &upload)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := UploadImageResponse{
		Img:   img.URL,
		Thumb: img.Thumb,
	}

	h.sendHTTPResponse(r.Context(), w, &result)
}

func (h *HTTPTransport) DeleteFile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	var req DeleteImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	if req.Filepath == "" {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	img := requests.Image{
		UserID:    userID,
		RequestID: id,
		URL:       req.Filepath,
	}

	if err := h.svc.DeleteImage(r.Context(), &img); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func (h *HTTPTransport) GuardList(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := parsePaginationRequest(r)
	if err != nil {
		h.log.Error("bad pagination request: %w", err)
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	req := requests.RequestListFilter{
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

	if err := validateFilterRequest(req); err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	res, count, err := h.svc.GuardRequestList(r.Context(), &req)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := RequestGuardListResponse{
		Data:  make([]*RequestForGuard, len(res)),
		Count: count,
	}

	for i := range res {
		result.Data[i] = &RequestForGuard{
			ID:          res[i].ID,
			UserID:      res[i].UserID,
			Type:        res[i].Type,
			Time:        res[i].Time,
			Description: res[i].Description,
			Status:      res[i].Status,
			UserName:    res[i].User.FirstName + " " + res[i].User.LastName,
			Phone:       res[i].User.Phone,
			Address:     res[i].User.Building.Name + ", " + res[i].User.Entry.Name,
			Apartment:   res[i].User.Apartment,
			Images:      res[i].ImagesURL,
		}
	}

	h.sendHTTPResponse(r.Context(), w, result)
}

func (h *HTTPTransport) GuardUpdateRequest(w http.ResponseWriter, r *http.Request) {
	var req GuardUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	if req.Status != "new" && req.Status != "closed" {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	data := requests.Request{
		ID:     id,
		Status: req.Status,
	}

	if err := h.svc.GuardUpdateRequest(r.Context(), &data); err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	h.sendHTTPResponse(r.Context(), w, nil)
}

func getIDFromQuery(r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errs.BadRequest
	}

	return uint(id), nil
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

func validateCreateRequest(r *RequestCreateRequest) error {
	// todo: add request status validation: oneof=all new closed
	reqTypes := map[string]struct{}{
		"taxi":     {},
		"guest":    {},
		"delivery": {},
		"noise":    {},
		"complain": {},
	}

	if _, ok := reqTypes[r.Type]; !ok {
		return errs.WrongRequestType
	}

	if r.Time <= 0 {
		return errs.WrongRequestDate
	}

	return nil
}

func validateFilterRequest(r requests.RequestListFilter) error {
	reqTypes := map[string]struct{}{
		"taxi":     {},
		"guest":    {},
		"delivery": {},
		"noise":    {},
		"complain": {},
		"all":      {},
	}

	if _, ok := reqTypes[r.Type]; !ok {
		return errs.WrongRequestType
	}

	if r.Place != "all" && r.Place != "kpp" {
		return errs.WrongRequestPlace
	}

	if r.Apartment != "" {
		if _, err := strconv.ParseFloat(r.Apartment, 64); err != nil {
			return errs.WrongApartment
		}
	}

	if r.Status != "all" && r.Status != "new" && r.Status != "closed" {
		return errs.WrongRequestStatus
	}

	return nil
}

func parsePaginationRequest(r *http.Request) (uint, uint, error) {
	_offset := r.URL.Query().Get("offset")
	_limit := r.URL.Query().Get("limit")

	if _offset == "" {
		return 0, 0, errs.EmptyOffset
	}

	if _limit == "" {
		return 0, 0, errs.EmptyLimit
	}

	offset, err := strconv.ParseUint(_offset, 10, 32)
	if err != nil {
		return 0, 0, errs.BadOffset
	}

	limit, err := strconv.ParseUint(_limit, 10, 32)
	if err != nil {
		return 0, 0, errs.BadLimit
	}

	if limit == 0 {
		return 0, 0, errs.LimitTooSmall
	}

	if limit > 200 {
		return 0, 0, errs.LimitTooBig
	}

	return uint(offset), uint(limit), nil
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
