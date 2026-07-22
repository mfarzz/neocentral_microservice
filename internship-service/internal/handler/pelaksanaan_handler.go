package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type PelaksanaanHandler struct {
	service service.PelaksanaanService
}

func NewPelaksanaanHandler(service service.PelaksanaanService) *PelaksanaanHandler {
	return &PelaksanaanHandler{service: service}
}

func (h *PelaksanaanHandler) GetInternshipHistory(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-User-Id")
	if studentID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	history, err := h.service.GetInternshipHistory(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

func (h *PelaksanaanHandler) GetLogbooks(w http.ResponseWriter, r *http.Request) {
	internshipID := chi.URLParam(r, "internshipId")
	
	logbooks, err := h.service.GetLogbooks(r.Context(), internshipID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logbooks)
}

func (h *PelaksanaanHandler) CreateLogbook(w http.ResponseWriter, r *http.Request) {
	internshipID := chi.URLParam(r, "internshipId")
	
	var logbook model.InternshipLogbook
	if err := json.NewDecoder(r.Body).Decode(&logbook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	logbook.InternshipID = internshipID

	if err := h.service.CreateLogbook(r.Context(), &logbook); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logbook)
}

func (h *PelaksanaanHandler) UpdateLogbook(w http.ResponseWriter, r *http.Request) {
	logbookID := chi.URLParam(r, "logbookId")
	
	var req struct {
		ActivityDescription string `json:"activityDescription"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateLogbook(r.Context(), logbookID, req.ActivityDescription); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logbook updated successfully"})
}

func (h *PelaksanaanHandler) UpdateInternshipDetails(w http.ResponseWriter, r *http.Request) {
	internshipID := chi.URLParam(r, "internshipId")
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateInternshipDetails(r.Context(), internshipID, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Internship details updated"})
}

// Other submissions (SubmitReport, UpdateCompletionCertificate, etc.)
func (h *PelaksanaanHandler) SubmitDocument(w http.ResponseWriter, r *http.Request) {
	internshipID := chi.URLParam(r, "internshipId")
	docType := chi.URLParam(r, "docType")
	
	var req struct {
		Title      string `json:"title"`
		DocumentID string `json:"documentId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var err error
	switch docType {
	case "report":
		err = h.service.SubmitReport(r.Context(), internshipID, req.Title, req.DocumentID)
	case "certificate":
		err = h.service.UpdateCompletionCertificate(r.Context(), internshipID, req.DocumentID)
	case "company-receipt":
		err = h.service.UpdateCompanyReceipt(r.Context(), internshipID, req.DocumentID)
	case "company-report":
		err = h.service.SubmitCompanyReport(r.Context(), internshipID, req.DocumentID)
	case "logbook-document":
		err = h.service.SubmitLogbookDocument(r.Context(), internshipID, req.DocumentID)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid document type"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document submitted successfully"})
}
