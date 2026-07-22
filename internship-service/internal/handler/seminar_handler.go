package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SeminarHandler struct {
	service service.SeminarService
}

func NewSeminarHandler(service service.SeminarService) *SeminarHandler {
	return &SeminarHandler{service: service}
}

func (h *SeminarHandler) GetUpcomingSeminars(w http.ResponseWriter, r *http.Request) {
	seminars, err := h.service.GetUpcoming(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(seminars)
}

func (h *SeminarHandler) RegisterSeminar(w http.ResponseWriter, r *http.Request) {
	var s model.InternshipSeminar
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	s.ID = uuid.NewString()

	if err := h.service.RegisterSeminar(r.Context(), &s); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *SeminarHandler) RegisterAudience(w http.ResponseWriter, r *http.Request) {
	seminarID := chi.URLParam(r, "id")
	studentID := r.Header.Get("X-User-Id")
	
	audience := model.InternshipSeminarAudience{
		ID:        uuid.NewString(),
		SeminarID: seminarID,
		StudentID: studentID,
	}

	if err := h.service.RegisterAudience(r.Context(), &audience); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Audience registered successfully"})
}

func (h *SeminarHandler) ValidateAudienceBulk(w http.ResponseWriter, r *http.Request) {
	seminarID := chi.URLParam(r, "id")
	lecturerID := r.Header.Get("X-User-Id")
	
	var req struct {
		StudentIDs []string `json:"studentIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.ValidateAudienceBulk(r.Context(), seminarID, req.StudentIDs, lecturerID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Audience validated successfully"})
}

func (h *SeminarHandler) ApproveSeminarBulk(w http.ResponseWriter, r *http.Request) {
	lecturerID := r.Header.Get("X-User-Id")
	
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.ApproveSeminarBulk(r.Context(), req.IDs, lecturerID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Seminars approved successfully"})
}
