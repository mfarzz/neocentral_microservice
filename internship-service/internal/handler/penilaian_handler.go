package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PenilaianHandler struct {
	service service.PenilaianService
}

func NewPenilaianHandler(service service.PenilaianService) *PenilaianHandler {
	return &PenilaianHandler{service: service}
}

func (h *PenilaianHandler) GetCPMKs(w http.ResponseWriter, r *http.Request) {
	academicYearID := r.URL.Query().Get("academicYearId")
	
	cpmks, err := h.service.GetCPMKs(r.Context(), academicYearID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cpmks)
}

func (h *PenilaianHandler) CreateCPMK(w http.ResponseWriter, r *http.Request) {
	var c model.InternshipCPMK
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	c.ID = uuid.NewString()

	if err := h.service.CreateCPMK(r.Context(), &c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *PenilaianHandler) CreateRubric(w http.ResponseWriter, r *http.Request) {
	var rub model.InternshipAssessmentRubric
	if err := json.NewDecoder(r.Body).Decode(&rub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	rub.ID = uuid.NewString()

	if err := h.service.CreateRubric(r.Context(), &rub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rub)
}

func (h *PenilaianHandler) SubmitLecturerAssessment(w http.ResponseWriter, r *http.Request) {
	internshipID := chi.URLParam(r, "internshipId")
	
	var req struct {
		Scores []service.ScoreInput `json:"scores"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.SubmitLecturerAssessment(r.Context(), internshipID, req.Scores); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Lecturer assessment saved successfully"})
}

func (h *PenilaianHandler) SubmitFieldAssessment(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	
	var req struct {
		Scores    []service.ScoreInput `json:"scores"`
		Notes     string               `json:"notes"`
		Signature string               `json:"signature"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.SubmitFieldAssessment(r.Context(), token, req.Scores, req.Notes, req.Signature); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Field assessment submitted successfully"})
}
