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

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.listByTask)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}

type commentCreateRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) create(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	var req commentCreateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Content == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "content is required"})
		return
	}

	user := r.Context().Value(model.UserKey).(*model.User)

	comment, err := h.commentService.Create(r.Context(), taskID, user.ID, req.Content)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to create comment"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) listByTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	comments, err := h.commentService.ListByTask(r.Context(), taskID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to list comments"})
		return
	}

	response.WriteJSON(w, http.StatusOK, comments)
}

type commentUpdateRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	var req commentUpdateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Content == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "content is required"})
		return
	}

	comment, err := h.commentService.Update(r.Context(), id, req.Content)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "comment not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to update comment"})
		return
	}

	response.WriteJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	err = h.commentService.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "comment not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to delete comment"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
