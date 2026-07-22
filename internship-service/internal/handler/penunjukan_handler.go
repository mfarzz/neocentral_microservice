package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type PenunjukanHandler struct {
	service service.PenunjukanService
}

func NewPenunjukanHandler(service service.PenunjukanService) *PenunjukanHandler {
	return &PenunjukanHandler{service: service}
}

func (h *PenunjukanHandler) BulkAssignSupervisor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InternshipIDs []string `json:"internshipIds"`
		SupervisorID  string   `json:"supervisorId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.BulkAssignSupervisor(r.Context(), req.InternshipIDs, req.SupervisorID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Supervisors successfully assigned to selected internships"})
}

func (h *PenunjukanHandler) GetSupervisorLetter(w http.ResponseWriter, r *http.Request) {
	supervisorID := chi.URLParam(r, "supervisorId")
	
	letter, err := h.service.GetSupervisorLetterDetail(r.Context(), supervisorID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Letter not found or error occurred"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(letter)
}

func (h *PenunjukanHandler) UpdateSupervisorLetter(w http.ResponseWriter, r *http.Request) {
	supervisorID := chi.URLParam(r, "supervisorId")
	
	var req model.InternshipSupervisorLetter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSupervisorLetter(r.Context(), supervisorID, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Supervisor letter successfully updated"})
}

func (h *PenunjukanHandler) GetLecturersWorkload(w http.ResponseWriter, r *http.Request) {
	academicYearID := r.URL.Query().Get("academicYearId")
	
	workload, err := h.service.GetLecturersWorkloadList(r.Context(), academicYearID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workload)
}
