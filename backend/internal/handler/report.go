package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/repository"
	"github.com/vshakitskiy/cis/internal/response"
)

type ReportHandler struct {
	reportRepo *repository.ReportRepo
}

func NewReportHandler(reportRepo *repository.ReportRepo) *ReportHandler {
	return &ReportHandler{reportRepo: reportRepo}
}

func (h *ReportHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.projectReport)

	return r
}

func (h *ReportHandler) projectReport(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid project id"})
		return
	}

	report, err := h.reportRepo.GetProjectReport(r.Context(), projectID)
	if err == pgx.ErrNoRows {
		response.WriteJSON(w, http.StatusNotFound, response.JSON{"error": "project not found"})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "failed to generate report"})
		return
	}

	response.WriteJSON(w, http.StatusOK, report)
}
