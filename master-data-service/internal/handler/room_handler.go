package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/pkg/response"
)

type RoomHandler struct {
	svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

// GetAll godoc
// @Summary      Get all rooms
// @Tags         Rooms
// @Produce      json
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/rooms [get]
// @Security     BearerAuth
func (h *RoomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Rooms retrieved", data)
}

// GetByID godoc
// @Summary      Get room by ID
// @Tags         Rooms
// @Produce      json
// @Param        id path string true "Room ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/rooms/{id} [get]
// @Security     BearerAuth
func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Room retrieved", data)
}

// Create godoc
// @Summary      Create a new room
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateRoomRequest true "Room data"
// @Success      201 {object} response.SuccessBody
// @Router       /master-data/rooms [post]
// @Security     BearerAuth
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Room created", data)
}

// Update godoc
// @Summary      Update a room
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param        id path string true "Room ID"
// @Param        body body dto.UpdateRoomRequest true "Updated fields"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/rooms/{id} [patch]
// @Security     BearerAuth
func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	data, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Room updated", data)
}

// Delete godoc
// @Summary      Delete a room
// @Tags         Rooms
// @Param        id path string true "Room ID"
// @Success      200 {object} response.SuccessBody
// @Router       /master-data/rooms/{id} [delete]
// @Security     BearerAuth
func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	response.Success(w, http.StatusOK, "Room deleted", nil)
}
