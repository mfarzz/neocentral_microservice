package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/pkg/apperror"
	"neocentral-go/pkg/response"
)

type AcademicYearHandler struct {
	svc *service.AcademicYearService
}

func NewAcademicYearHandler(svc *service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{svc: svc}
}

// GetAll godoc
// @Summary      Get all academic years
// @Tags         Academic Years
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/academic-years [get]
// @Security     BearerAuth
func (h *AcademicYearHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Academic years retrieved", data)
}

// GetActive godoc
// @Summary      Get the currently active academic year
// @Tags         Academic Years
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/academic-years/active [get]
// @Security     BearerAuth
func (h *AcademicYearHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetActive(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Active academic year retrieved", data)
}

// GetByID godoc
// @Summary      Get academic year by ID
// @Tags         Academic Years
// @Produce      json
// @Param        id path string true "Academic Year ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/academic-years/{id} [get]
// @Security     BearerAuth
func (h *AcademicYearHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Academic year retrieved", data)
}

// Create godoc
// @Summary      Create a new academic year
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateAcademicYearRequest true "Academic Year data"
// @Success      201 {object} response.SuccessBody
// @Router       /master-data/academic-years [post]
// @Security     BearerAuth
func (h *AcademicYearHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAcademicYearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Academic year created", data)
}

// Update godoc
// @Summary      Update an academic year
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        id path string true "Academic Year ID"
// @Param        body body dto.UpdateAcademicYearRequest true "Updated fields"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/academic-years/{id} [patch]
// @Security     BearerAuth
func (h *AcademicYearHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateAcademicYearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Academic year updated", data)
}

// Delete godoc
// @Summary      Delete an academic year
// @Tags         Academic Years
// @Param        id path string true "Academic Year ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/academic-years/{id} [delete]
// @Security     BearerAuth
func (h *AcademicYearHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Academic year deleted", nil)
}

// handleError maps apperror types to HTTP responses.
func handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}
	response.Error(w, http.StatusInternalServerError, "Internal server error")
}
