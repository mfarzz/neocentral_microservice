package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/model"
	"neocentral-go/master-data-service/internal/repository"
	"neocentral-go/pkg/apperror"
)

type InternshipHolidayService struct {
	repo repository.InternshipHolidayRepository
}

func NewInternshipHolidayService(repo repository.InternshipHolidayRepository) *InternshipHolidayService {
	return &InternshipHolidayService{repo: repo}
}

func (s *InternshipHolidayService) GetAll(ctx context.Context) ([]dto.InternshipHolidayResponse, error) {
	holidays, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("Failed to fetch internship holidays", err)
	}

	var responses []dto.InternshipHolidayResponse
	for _, h := range holidays {
		responses = append(responses, buildInternshipHolidayResponse(&h))
	}

	if responses == nil {
		responses = make([]dto.InternshipHolidayResponse, 0)
	}

	return responses, nil
}

func (s *InternshipHolidayService) Create(ctx context.Context, req dto.CreateInternshipHolidayRequest) (*dto.InternshipHolidayResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.HolidayDate)
	if err != nil {
		return nil, apperror.BadRequest("Invalid holidayDate format. Expected YYYY-MM-DD")
	}

	// Check if already exists
	existing, err := s.repo.FindByDate(ctx, req.HolidayDate)
	if err != nil {
		return nil, apperror.InternalWrap("Database error", err)
	}
	if existing != nil {
		return nil, apperror.Conflict("Internship holiday with this date already exists")
	}

	holiday := &model.InternshipHoliday{
		ID:          uuid.New().String(),
		HolidayDate: parsedDate,
		Name:        req.Name,
	}

	if err := s.repo.Create(ctx, holiday); err != nil {
		return nil, apperror.InternalWrap("Failed to create internship holiday", err)
	}

	resp := buildInternshipHolidayResponse(holiday)
	return &resp, nil
}

func (s *InternshipHolidayService) Update(ctx context.Context, id string, req dto.UpdateInternshipHolidayRequest) (*dto.InternshipHolidayResponse, error) {
	holiday, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("Database error", err)
	}
	if holiday == nil {
		return nil, apperror.NotFound("Internship holiday not found")
	}

	if req.HolidayDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.HolidayDate)
		if err != nil {
			return nil, apperror.BadRequest("Invalid holidayDate format. Expected YYYY-MM-DD")
		}
		// check conflict
		existing, err := s.repo.FindByDate(ctx, *req.HolidayDate)
		if err != nil {
			return nil, apperror.InternalWrap("Database error", err)
		}
		if existing != nil && existing.ID != id {
			return nil, apperror.Conflict("Internship holiday with this date already exists")
		}
		holiday.HolidayDate = parsedDate
	}
	if req.Name != nil {
		holiday.Name = req.Name
	}

	if err := s.repo.Update(ctx, holiday); err != nil {
		return nil, apperror.InternalWrap("Failed to update internship holiday", err)
	}

	resp := buildInternshipHolidayResponse(holiday)
	return &resp, nil
}

func (s *InternshipHolidayService) Delete(ctx context.Context, id string) error {
	holiday, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.InternalWrap("Database error", err)
	}
	if holiday == nil {
		return apperror.NotFound("Internship holiday not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperror.InternalWrap("Failed to delete internship holiday", err)
	}

	return nil
}

func buildInternshipHolidayResponse(h *model.InternshipHoliday) dto.InternshipHolidayResponse {
	return dto.InternshipHolidayResponse{
		ID:          h.ID,
		HolidayDate: h.HolidayDate.Format("2006-01-02"),
		Name:        h.Name,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}
