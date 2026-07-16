package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/pkg/response"
)

type ThesisHandler struct {
	svc *service.ThesisTopicService
}

func NewThesisHandler(svc *service.ThesisTopicService) *ThesisHandler {
	return &ThesisHandler{svc: svc}
}

// ── Topics ───────────────────────────────────

// GetAllTopics godoc
// @Summary      Get all thesis topics
// @Tags         Thesis Topics
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/thesis-topics [get]
// @Security     BearerAuth
func (h *ThesisHandler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAllTopics(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Thesis topics retrieved", data)
}

// GetTopicByID godoc
// @Summary      Get thesis topic by ID
// @Tags         Thesis Topics
// @Produce      json
// @Param        id path string true "Thesis Topic ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/thesis-topics/{id} [get]
// @Security     BearerAuth
func (h *ThesisHandler) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := h.svc.GetTopicByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Thesis topic retrieved", data)
}

// CreateTopic godoc
// @Summary      Create a new thesis topic
// @Tags         Thesis Topics
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateThesisTopicRequest true "Topic data"
// @Success      201 {object} response.SuccessBody
// @Router       /master-data/thesis-topics [post]
// @Security     BearerAuth
func (h *ThesisHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateThesisTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.CreateTopic(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Thesis topic created", data)
}

// UpdateTopic godoc
// @Summary      Update a thesis topic
// @Tags         Thesis Topics
// @Accept       json
// @Produce      json
// @Param        id path string true "Thesis Topic ID"
// @Param        body body dto.UpdateThesisTopicRequest true "Updated fields"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/thesis-topics/{id} [patch]
// @Security     BearerAuth
func (h *ThesisHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateThesisTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.UpdateTopic(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Thesis topic updated", data)
}

// DeleteTopic godoc
// @Summary      Delete a thesis topic
// @Tags         Thesis Topics
// @Param        id path string true "Thesis Topic ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/thesis-topics/{id} [delete]
// @Security     BearerAuth
func (h *ThesisHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteTopic(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Thesis topic deleted", nil)
}

// ── Statuses ─────────────────────────────────

// GetAllStatuses godoc
// @Summary      Get all thesis statuses
// @Tags         Thesis Statuses
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/thesis-statuses [get]
// @Security     BearerAuth
func (h *ThesisHandler) GetAllStatuses(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAllStatuses(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Thesis statuses retrieved", data)
}
