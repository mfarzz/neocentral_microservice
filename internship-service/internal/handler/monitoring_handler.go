package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type MonitoringHandler struct {
	service service.MonitoringService
}

func NewMonitoringHandler(service service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{service: service}
}

func (h *MonitoringHandler) GetInternshipList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := map[string]interface{}{} // extract from query params normally

	data, total, err := h.service.GetInternshipList(r.Context(), filters, page, pageSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"total":   total,
	})
}

func (h *MonitoringHandler) GetInternshipDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := h.service.GetInternshipDetail(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *MonitoringHandler) GetGradeRecap(w http.ResponseWriter, r *http.Request) {
	academicYear := r.URL.Query().Get("academicYearId")
	data, err := h.service.GetGradeRecap(r.Context(), academicYear)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *MonitoringHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		DocumentType string `json:"documentType"`
		Status       string `json:"status"`
		Notes        string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.VerifyDocument(r.Context(), id, req.DocumentType, req.Status, req.Notes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document verified"})
}

func (h *MonitoringHandler) RejectFinalReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.RejectFinalReport(r.Context(), id, req.Notes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Final report rejected"})
}

func (h *MonitoringHandler) GetMonitoringStats(w http.ResponseWriter, r *http.Request) {
	academicYear := r.URL.Query().Get("academicYearId")
	data, err := h.service.GetMonitoringStats(r.Context(), academicYear)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
