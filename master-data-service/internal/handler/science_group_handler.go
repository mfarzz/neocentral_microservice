package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/pkg/response"
)

type ScienceGroupHandler struct {
	svc *service.ScienceGroupService
}

func NewScienceGroupHandler(svc *service.ScienceGroupService) *ScienceGroupHandler {
	return &ScienceGroupHandler{svc: svc}
}

// GetAll godoc
// @Summary      Get all science groups
// @Tags         Science Groups
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/science-groups [get]
// @Security     BearerAuth
func (h *ScienceGroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Science groups retrieved", data)
}

// GetByID godoc
// @Summary      Get science group by ID
// @Tags         Science Groups
// @Produce      json
// @Param        id path string true "Science Group ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/science-groups/{id} [get]
// @Security     BearerAuth
func (h *ScienceGroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Science group retrieved", data)
}

// Create godoc
// @Summary      Create a new science group
// @Tags         Science Groups
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateScienceGroupRequest true "Science Group data"
// @Success      201 {object} response.SuccessBody
// @Router       /master-data/science-groups [post]
// @Security     BearerAuth
func (h *ScienceGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateScienceGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Science group created", data)
}

// Update godoc
// @Summary      Update a science group
// @Tags         Science Groups
// @Accept       json
// @Produce      json
// @Param        id path string true "Science Group ID"
// @Param        body body dto.UpdateScienceGroupRequest true "Updated fields"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/science-groups/{id} [patch]
// @Security     BearerAuth
func (h *ScienceGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateScienceGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Science group updated", data)
}

// Delete godoc
// @Summary      Delete a science group
// @Tags         Science Groups
// @Param        id path string true "Science Group ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/science-groups/{id} [delete]
// @Security     BearerAuth
func (h *ScienceGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Science group deleted", nil)
}
