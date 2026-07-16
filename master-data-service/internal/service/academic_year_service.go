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

type AcademicYearService struct {
	repo repository.AcademicYearRepository
}

func NewAcademicYearService(repo repository.AcademicYearRepository) *AcademicYearService {
	return &AcademicYearService{repo: repo}
}

func (s *AcademicYearService) GetAll(ctx context.Context) ([]dto.AcademicYearResponse, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("failed to fetch academic years", err)
	}
	result := make([]dto.AcademicYearResponse, len(list))
	for i, ay := range list {
		result[i] = toAcademicYearResponse(ay)
	}
	return result, nil
}

func (s *AcademicYearService) GetByID(ctx context.Context, id string) (*dto.AcademicYearResponse, error) {
	ay, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if ay == nil {
		return nil, apperror.NotFound("Academic year not found")
	}
	resp := toAcademicYearResponse(*ay)
	return &resp, nil
}

func (s *AcademicYearService) GetActive(ctx context.Context) (*dto.AcademicYearResponse, error) {
	ay, err := s.repo.FindActive(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if ay == nil {
		return nil, apperror.NotFound("No active academic year found")
	}
	resp := toAcademicYearResponse(*ay)
	return &resp, nil
}

func (s *AcademicYearService) Create(ctx context.Context, req dto.CreateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	ay := model.AcademicYear{
		ID:       uuid.NewString(),
		Semester: req.Semester,
		IsActive: req.IsActive,
	}
	if req.Year != "" {
		ay.Year = &req.Year
	}
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, apperror.BadRequest("Invalid startDate format, use YYYY-MM-DD")
		}
		ay.StartDate = &t
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, apperror.BadRequest("Invalid endDate format, use YYYY-MM-DD")
		}
		ay.EndDate = &t
	}

	// If setting as active, deactivate all others first
	if req.IsActive {
		if err := s.repo.DeactivateAll(ctx); err != nil {
			return nil, apperror.InternalWrap("failed to deactivate previous academic years", err)
		}
	}

	if err := s.repo.Create(ctx, &ay); err != nil {
		return nil, apperror.InternalWrap("failed to create academic year", err)
	}
	resp := toAcademicYearResponse(ay)
	return &resp, nil
}

func (s *AcademicYearService) Update(ctx context.Context, id string, req dto.UpdateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	ay, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if ay == nil {
		return nil, apperror.NotFound("Academic year not found")
	}

	if req.Semester != nil {
		ay.Semester = *req.Semester
	}
	if req.Year != nil {
		ay.Year = req.Year
	}
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, apperror.BadRequest("Invalid startDate format")
		}
		ay.StartDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, apperror.BadRequest("Invalid endDate format")
		}
		ay.EndDate = &t
	}
	if req.IsActive != nil {
		if *req.IsActive {
			_ = s.repo.DeactivateAll(ctx)
		}
		ay.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, ay); err != nil {
		return nil, apperror.InternalWrap("failed to update academic year", err)
	}
	resp := toAcademicYearResponse(*ay)
	return &resp, nil
}

func (s *AcademicYearService) Delete(ctx context.Context, id string) error {
	ay, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.InternalWrap("database error", err)
	}
	if ay == nil {
		return apperror.NotFound("Academic year not found")
	}
	return s.repo.Delete(ctx, id)
}

func toAcademicYearResponse(ay model.AcademicYear) dto.AcademicYearResponse {
	return dto.AcademicYearResponse{
		ID:        ay.ID,
		Semester:  ay.Semester,
		Year:      ay.Year,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}
}
