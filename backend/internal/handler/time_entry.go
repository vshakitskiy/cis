package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/response"
	"github.com/vshakitskiy/cis/internal/service"
)

type TimeEntryHandler struct {
	timeEntryService *service.TimeEntryService
}

func NewTimeEntryHandler(timeEntryService *service.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{timeEntryService: timeEntryService}
}

func (h *TimeEntryHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.listByTask)
	r.Delete("/{id}", h.delete)

	return r
}

type timeEntryCreateRequest struct {
	Minutes     int        `json:"minutes"`
	Description *string    `json:"description"`
	Date        model.Date `json:"date"`
}

func (h *TimeEntryHandler) create(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	var req timeEntryCreateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Minutes <= 0 {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "minutes must be positive"})
		return
	}

	user := r.Context().Value(model.UserKey).(*model.User)

	entry, err := h.timeEntryService.Create(r.Context(), taskID, user.ID, req.Minutes, req.Description, req.Date)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to create time entry"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, entry)
}

func (h *TimeEntryHandler) listByTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	entries, err := h.timeEntryService.ListByTask(r.Context(), taskID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to list time entries"})
		return
	}

	response.WriteJSON(w, http.StatusOK, entries)
}

func (h *TimeEntryHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	err = h.timeEntryService.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "time entry not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to delete time entry"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
