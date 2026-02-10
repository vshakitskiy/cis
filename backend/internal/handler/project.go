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

type ProjectHandler struct {
	projectService *service.ProjectService
	taskHandler    *TaskHandler
	reportHandler  *ReportHandler
}

func NewProjectHandler(projectService *service.ProjectService, taskHandler *TaskHandler, reportHandler *ReportHandler) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		taskHandler:    taskHandler,
		reportHandler:  reportHandler,
	}
}

func (h *ProjectHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	r.Mount("/{projectID}/tasks", h.taskHandler.Routes())
	r.Mount("/{projectID}/report", h.reportHandler.Routes())

	return r
}

type projectRequest struct {
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	Deadline    *model.Date `json:"deadline"`
}

func (h *ProjectHandler) create(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "name is required"})
		return
	}

	user := r.Context().Value(model.UserKey).(*model.User)

	project, err := h.projectService.Create(r.Context(), req.Name, req.Description, user.ID, req.Deadline)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to create project"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.List(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to list projects"})
		return
	}

	response.WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	project, err := h.projectService.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "project not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to get project"})
		return
	}

	response.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	var req projectRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "name is required"})
		return
	}

	project, err := h.projectService.Update(r.Context(), id, req.Name, req.Description, req.Deadline)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "project not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to update project"})
		return
	}

	response.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	err = h.projectService.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "project not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to delete project"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
