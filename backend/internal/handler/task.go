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

type TaskHandler struct {
	taskService       *service.TaskService
	timeEntryHandler  *TimeEntryHandler
	commentHandler    *CommentHandler
	attachmentHandler *AttachmentHandler
}

func NewTaskHandler(taskService *service.TaskService, timeEntryHandler *TimeEntryHandler, commentHandler *CommentHandler, attachmentHandler *AttachmentHandler) *TaskHandler {
	return &TaskHandler{
		taskService:       taskService,
		timeEntryHandler:  timeEntryHandler,
		commentHandler:    commentHandler,
		attachmentHandler: attachmentHandler,
	}
}

func (h *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.listByProject)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	r.Mount("/{taskID}/time-entries", h.timeEntryHandler.Routes())
	r.Mount("/{taskID}/comments", h.commentHandler.Routes())
	r.Mount("/{taskID}/attachments", h.attachmentHandler.Routes())

	return r
}

type taskCreateRequest struct {
	AssigneeID  *int64             `json:"assignee_id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Priority    model.TaskPriority `json:"priority"`
	Deadline    *model.Date        `json:"deadline"`
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid project id"})
		return
	}

	var req taskCreateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "title is required"})
		return
	}

	if req.Priority == "" {
		req.Priority = model.PriorityMedium
	}

	task, err := h.taskService.Create(r.Context(), projectID, req.AssigneeID, req.Title, req.Description, req.Priority, req.Deadline)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "project not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to create task"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) listByProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid project id"})
		return
	}

	tasks, err := h.taskService.ListByProject(r.Context(), projectID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to list tasks"})
		return
	}

	response.WriteJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	task, err := h.taskService.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to get task"})
		return
	}

	response.WriteJSON(w, http.StatusOK, task)
}

type taskUpdateRequest struct {
	AssigneeID  *int64             `json:"assignee_id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Status      model.TaskStatus   `json:"status"`
	Priority    model.TaskPriority `json:"priority"`
	Deadline    *model.Date        `json:"deadline"`
}

func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	var req taskUpdateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Title == "" || req.Status == "" || req.Priority == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "title, status, and priority are required"})
		return
	}

	task, err := h.taskService.Update(r.Context(), id, req.AssigneeID, req.Title, req.Description, req.Status, req.Priority, req.Deadline)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to update task"})
		return
	}

	response.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	err = h.taskService.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to delete task"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
