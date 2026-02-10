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

const maxUploadSize = 10 << 20 // 10 MB

type AttachmentHandler struct {
	attachmentService *service.AttachmentService
}

func NewAttachmentHandler(attachmentService *service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{attachmentService: attachmentService}
}

func (h *AttachmentHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.upload)
	r.Get("/", h.listByTask)
	r.Get("/{id}/download", h.download)
	r.Delete("/{id}", h.delete)

	return r
}

func (h *AttachmentHandler) upload(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "file too large (max 10MB)"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "file is required"})
		return
	}
	defer file.Close()

	user := r.Context().Value(model.UserKey).(*model.User)

	attachment, err := h.attachmentService.Create(r.Context(), taskID, user.ID, header.Filename, header.Size, file)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "task not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to upload file"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, attachment)
}

func (h *AttachmentHandler) listByTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid task id"})
		return
	}

	attachments, err := h.attachmentService.ListByTask(r.Context(), taskID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to list attachments"})
		return
	}

	response.WriteJSON(w, http.StatusOK, attachments)
}

func (h *AttachmentHandler) download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	attachment, err := h.attachmentService.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "attachment not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to get attachment"})
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+attachment.Filename+"\"")
	http.ServeFile(w, r, attachment.Filepath)
}

func (h *AttachmentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid id"})
		return
	}

	err = h.attachmentService.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "attachment not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to delete attachment"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
