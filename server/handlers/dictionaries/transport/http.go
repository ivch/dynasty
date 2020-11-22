package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/dictionaries"
)

type DictionaryService interface {
	EntriesList(ctx context.Context, buildingID uint) ([]*dictionaries.Entry, error)
	BuildingsList(ctx context.Context) ([]*dictionaries.Building, error)
}

type HTTPTransport struct {
	svc    DictionaryService
	log    logger.Logger
	router chi.Router
}

func (h *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// NewHTTPTransport returns a new instance of HTTPTransport.
func NewHTTPTransport(log logger.Logger, svc DictionaryService, mdl ...func(http.Handler) http.Handler) http.Handler {
	h := &HTTPTransport{log: log, router: chi.NewRouter().With(mdl...), svc: svc}
	h.attachRoutes()
	return h
}

func (h *HTTPTransport) attachRoutes() {
	h.router.Get("/v1/buildings", h.Buildings)
	h.router.Get("/v1/building/{id}/entries", h.Entries)
}

func (h *HTTPTransport) Buildings(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.BuildingsList(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]*Building, len(res))
	for i := range res {
		result[i] = &Building{
			ID:      res[i].ID,
			Name:    res[i].Name,
			Address: res[i].Address,
		}
	}

	h.sendHTTPResponse(r.Context(), w, BuildingsDictionaryResposnse{Data: result})
}

func (h *HTTPTransport) Entries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		h.sendError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	res, err := h.svc.EntriesList(r.Context(), uint(id))
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]*Entry, len(res))
	for i := range res {
		result[i] = &Entry{
			ID:   res[i].ID,
			Name: res[i].Name,
		}
	}

	h.sendHTTPResponse(r.Context(), w, EntriesDictionaryResponse{Data: result})
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
