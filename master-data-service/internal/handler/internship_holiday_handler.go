package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/pkg/response"
)

type InternshipHolidayHandler struct {
	svc *service.InternshipHolidayService
}

func NewInternshipHolidayHandler(svc *service.InternshipHolidayService) *InternshipHolidayHandler {
	return &InternshipHolidayHandler{svc: svc}
}

// @Summary      Get all internship holidays
// @Description  Get a list of all internship holidays
// @Tags         Internship Holidays
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]dto.InternshipHolidayResponse}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /internship-holidays [get]
func (h *InternshipHolidayHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Success", resp)
}

// @Summary      Create internship holiday
// @Description  Create a new internship holiday
// @Tags         Internship Holidays
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateInternshipHolidayRequest true "Request Body"
// @Success      201  {object}  response.APIResponse{data=dto.InternshipHolidayResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      409  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /internship-holidays [post]
func (h *InternshipHolidayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInternshipHolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Internship holiday created successfully", resp)
}

// @Summary      Update internship holiday
// @Description  Update an existing internship holiday
// @Tags         Internship Holidays
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Holiday ID"
// @Param        request body dto.UpdateInternshipHolidayRequest true "Request Body"
// @Success      200  {object}  response.APIResponse{data=dto.InternshipHolidayResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /internship-holidays/{id} [put]
func (h *InternshipHolidayHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateInternshipHolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Internship holiday updated successfully", resp)
}

// @Summary      Delete internship holiday
// @Description  Delete an existing internship holiday
// @Tags         Internship Holidays
// @Produce      json
// @Param        id   path      string  true  "Holiday ID"
// @Success      200  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /internship-holidays/{id} [delete]
func (h *InternshipHolidayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Internship holiday deleted successfully", nil)
}
