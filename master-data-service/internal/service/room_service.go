package service

import (
	"context"

	"github.com/google/uuid"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/model"
	"neocentral-go/master-data-service/internal/repository"
	"neocentral-go/pkg/apperror"
)

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) GetAll(ctx context.Context) ([]dto.RoomResponse, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("failed to fetch rooms", err)
	}
	result := make([]dto.RoomResponse, len(list))
	for i, r := range list {
		result[i] = toRoomResponse(r)
	}
	return result, nil
}

func (s *RoomService) GetByID(ctx context.Context, id string) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if room == nil {
		return nil, apperror.NotFound("Room not found")
	}
	resp := toRoomResponse(*room)
	return &resp, nil
}

func (s *RoomService) Create(ctx context.Context, req dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	room := model.Room{
		ID:   uuid.NewString(),
		Name: req.Name,
	}
	if req.Location != "" {
		room.Location = &req.Location
	}
	if req.Capacity != nil {
		room.Capacity = req.Capacity
	}

	if err := s.repo.Create(ctx, &room); err != nil {
		return nil, apperror.InternalWrap("failed to create room", err)
	}
	resp := toRoomResponse(room)
	return &resp, nil
}

func (s *RoomService) Update(ctx context.Context, id string, req dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if room == nil {
		return nil, apperror.NotFound("Room not found")
	}

	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.Location != nil {
		room.Location = req.Location
	}
	if req.Capacity != nil {
		room.Capacity = req.Capacity
	}

	if err := s.repo.Update(ctx, room); err != nil {
		return nil, apperror.InternalWrap("failed to update room", err)
	}
	resp := toRoomResponse(*room)
	return &resp, nil
}

func (s *RoomService) Delete(ctx context.Context, id string) error {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.InternalWrap("database error", err)
	}
	if room == nil {
		return apperror.NotFound("Room not found")
	}
	return s.repo.Delete(ctx, id)
}

func toRoomResponse(r model.Room) dto.RoomResponse {
	return dto.RoomResponse{
		ID:        r.ID,
		Name:      r.Name,
		Location:  r.Location,
		Capacity:  r.Capacity,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
