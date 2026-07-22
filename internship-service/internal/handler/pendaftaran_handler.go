package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type PendaftaranHandler struct {
	service service.PendaftaranService
}

func NewPendaftaranHandler(service service.PendaftaranService) *PendaftaranHandler {
	return &PendaftaranHandler{service: service}
}

// CreateProposal creates a new internship proposal.
// @Summary      Create Proposal
// @Description  Creates a new internship proposal
// @Tags         pendaftaran
// @Accept       json
// @Produce      json
// @Param        proposal body model.InternshipProposal true "Proposal Request"
// @Success      201  {object}  model.InternshipProposal
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /pendaftaran/proposals [post]
func (h *PendaftaranHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var proposal model.InternshipProposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.CreateProposal(r.Context(), &proposal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proposal)
}

// GetProposals retrieves all internship proposals.
// @Summary      Get Proposals
// @Description  Retrieves all internship proposals
// @Tags         pendaftaran
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.InternshipProposal
// @Failure      500  {object}  map[string]string
// @Router       /pendaftaran/proposals [get]
func (h *PendaftaranHandler) GetProposals(w http.ResponseWriter, r *http.Request) {
	proposals, err := h.service.GetProposals(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposals)
}

// GetProposalByID retrieves an internship proposal by ID.
// @Summary      Get Proposal by ID
// @Description  Retrieves an internship proposal by its ID
// @Tags         pendaftaran
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Proposal ID"
// @Success      200  {object}  model.InternshipProposal
// @Failure      404  {object}  map[string]string
// @Router       /pendaftaran/proposals/{id} [get]
func (h *PendaftaranHandler) GetProposalByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	proposal, err := h.service.GetProposalByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "proposal not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

func (h *PendaftaranHandler) ListStudentProposals(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-User-Id")
	if studentID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "X-User-Id header is missing"})
		return
	}

	proposals, err := h.service.ListStudentProposals(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposals)
}

// UpdateProposal updates an existing internship proposal.
// @Summary      Update Proposal
// @Description  Updates an internship proposal by ID
// @Tags         pendaftaran
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Proposal ID"
// @Param        proposal body model.InternshipProposal true "Update Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /pendaftaran/proposals/{id} [put]
func (h *PendaftaranHandler) UpdateProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var proposal model.InternshipProposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProposal(r.Context(), id, &proposal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Proposal successfully updated"})
}

func (h *PendaftaranHandler) DeleteProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteProposal(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Proposal successfully deleted"})
}

func (h *PendaftaranHandler) UpdateProposalStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProposalStatus(r.Context(), id, model.ProposalStatus(req.Status)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Proposal status successfully updated"})
}

func (h *PendaftaranHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.ListCompanies(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}
