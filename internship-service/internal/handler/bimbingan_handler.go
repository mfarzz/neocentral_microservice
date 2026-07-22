package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/google/uuid"
)

type BimbinganHandler struct {
	service service.BimbinganService
}

func NewBimbinganHandler(service service.BimbinganService) *BimbinganHandler {
	return &BimbinganHandler{service: service}
}

func (h *BimbinganHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	academicYearID := r.URL.Query().Get("academicYearId")
	
	questions, err := h.service.GetQuestions(r.Context(), academicYearID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func (h *BimbinganHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var q model.InternshipGuidanceQuestion
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	q.ID = uuid.NewString()

	if err := h.service.CreateQuestion(r.Context(), &q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func (h *BimbinganHandler) GetLecturerCriteria(w http.ResponseWriter, r *http.Request) {
	academicYearID := r.URL.Query().Get("academicYearId")
	
	criteria, err := h.service.GetLecturerCriteria(r.Context(), academicYearID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(criteria)
}

func (h *BimbinganHandler) CreateLecturerCriteria(w http.ResponseWriter, r *http.Request) {
	var c model.InternshipGuidanceLecturerCriteria
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	c.ID = uuid.NewString()

	if err := h.service.CreateLecturerCriteria(r.Context(), &c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *BimbinganHandler) GetStudentGuidance(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-User-Id")
	
	guidance, err := h.service.GetStudentGuidanceSessions(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guidance)
}

func (h *BimbinganHandler) SubmitStudentGuidance(w http.ResponseWriter, r *http.Request) {
	var session model.InternshipGuidanceSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	session.ID = uuid.NewString()

	if err := h.service.SubmitGuidanceSession(r.Context(), &session); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (h *BimbinganHandler) SubmitLecturerEvaluation(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	var answers []model.InternshipGuidanceLecturerAnswer
	if err := json.NewDecoder(r.Body).Decode(&answers); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	if err := h.service.SubmitLecturerEvaluation(r.Context(), sessionID, answers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
