package service

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type PenunjukanService interface {
	BulkAssignSupervisor(ctx context.Context, internshipIDs []string, supervisorID string) error
	GetSupervisorLetterDetail(ctx context.Context, supervisorID string) (*model.InternshipSupervisorLetter, error)
	UpdateSupervisorLetter(ctx context.Context, supervisorID string, data *model.InternshipSupervisorLetter) error
	GetLecturersWorkloadList(ctx context.Context, academicYearID string) ([]map[string]interface{}, error)
}

type penunjukanService struct {
	internshipRepo repository.InternshipRepository
	letterRepo     repository.SupervisorLetterRepository
}

func NewPenunjukanService(internshipRepo repository.InternshipRepository, letterRepo repository.SupervisorLetterRepository) PenunjukanService {
	return &penunjukanService{
		internshipRepo: internshipRepo,
		letterRepo:     letterRepo,
	}
}

func (s *penunjukanService) BulkAssignSupervisor(ctx context.Context, internshipIDs []string, supervisorID string) error {
	return s.internshipRepo.UpdateSupervisorBulk(ctx, internshipIDs, supervisorID)
}

func (s *penunjukanService) GetSupervisorLetterDetail(ctx context.Context, supervisorID string) (*model.InternshipSupervisorLetter, error) {
	return s.letterRepo.GetBySupervisorID(ctx, supervisorID)
}

func (s *penunjukanService) UpdateSupervisorLetter(ctx context.Context, supervisorID string, data *model.InternshipSupervisorLetter) error {
	existing, err := s.letterRepo.GetBySupervisorID(ctx, supervisorID)
	if err != nil {
		// If not found, create new
		data.SupervisorID = supervisorID
		return s.letterRepo.Create(ctx, data)
	}

	// Update existing
	existing.DocumentNumber = data.DocumentNumber
	existing.DateIssued = data.DateIssued
	existing.StartDate = data.StartDate
	existing.EndDate = data.EndDate
	existing.DocumentID = data.DocumentID
	existing.SignedByID = data.SignedByID
	existing.SignedAsRoleID = data.SignedAsRoleID
	
	return s.letterRepo.Update(ctx, existing)
}

func (s *penunjukanService) GetLecturersWorkloadList(ctx context.Context, academicYearID string) ([]map[string]interface{}, error) {
	return s.internshipRepo.GetLecturerWorkload(ctx, academicYearID)
}
