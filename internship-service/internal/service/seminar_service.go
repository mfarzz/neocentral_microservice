package service

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type SeminarService interface {
	GetUpcoming(ctx context.Context) ([]model.InternshipSeminar, error)
	GetDetail(ctx context.Context, id string) (*model.InternshipSeminar, error)
	RegisterSeminar(ctx context.Context, seminar *model.InternshipSeminar) error
	UpdateSeminar(ctx context.Context, seminar *model.InternshipSeminar) error
	
	// Lecturer approvals
	ApproveSeminarBulk(ctx context.Context, ids []string, lecturerID string) error
	RejectSeminar(ctx context.Context, id string, notes string, lecturerID string) error
	CompleteSeminar(ctx context.Context, id string, lecturerID string) error
	
	// Audience
	RegisterAudience(ctx context.Context, audience *model.InternshipSeminarAudience) error
	UnregisterAudience(ctx context.Context, id string, studentID string) error
	ValidateAudienceBulk(ctx context.Context, id string, studentIDs []string, lecturerID string) error
	UnvalidateAudience(ctx context.Context, id string, studentID string, lecturerID string) error
}

type seminarService struct {
	seminarRepo repository.SeminarRepository
}

func NewSeminarService(seminarRepo repository.SeminarRepository) SeminarService {
	return &seminarService{
		seminarRepo: seminarRepo,
	}
}

func (s *seminarService) GetUpcoming(ctx context.Context) ([]model.InternshipSeminar, error) {
	return s.seminarRepo.GetUpcoming(ctx)
}

func (s *seminarService) GetDetail(ctx context.Context, id string) (*model.InternshipSeminar, error) {
	return s.seminarRepo.GetByID(ctx, id)
}

func (s *seminarService) RegisterSeminar(ctx context.Context, seminar *model.InternshipSeminar) error {
	return s.seminarRepo.Create(ctx, seminar)
}

func (s *seminarService) UpdateSeminar(ctx context.Context, seminar *model.InternshipSeminar) error {
	return s.seminarRepo.Update(ctx, seminar)
}

func (s *seminarService) ApproveSeminarBulk(ctx context.Context, ids []string, lecturerID string) error {
	for _, id := range ids {
		seminar, err := s.seminarRepo.GetByID(ctx, id)
		if err != nil {
			continue // skip or return err based on requirements
		}
		seminar.Status = model.SeminarApproved
		seminar.ApprovedBy = lecturerID
		s.seminarRepo.Update(ctx, seminar)
	}
	return nil
}

func (s *seminarService) RejectSeminar(ctx context.Context, id string, notes string, lecturerID string) error {
	seminar, err := s.seminarRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	seminar.Status = model.SeminarRejected
	seminar.SupervisorNotes = notes
	seminar.ApprovedBy = lecturerID
	return s.seminarRepo.Update(ctx, seminar)
}

func (s *seminarService) CompleteSeminar(ctx context.Context, id string, lecturerID string) error {
	seminar, err := s.seminarRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	seminar.Status = model.SeminarCompleted
	// additional logic e.g., verified by lecturerID
	return s.seminarRepo.Update(ctx, seminar)
}

func (s *seminarService) RegisterAudience(ctx context.Context, audience *model.InternshipSeminarAudience) error {
	return s.seminarRepo.AddAudience(ctx, audience)
}

func (s *seminarService) UnregisterAudience(ctx context.Context, id string, studentID string) error {
	return s.seminarRepo.RemoveAudience(ctx, id, studentID)
}

func (s *seminarService) ValidateAudienceBulk(ctx context.Context, id string, studentIDs []string, lecturerID string) error {
	return s.seminarRepo.ValidateAudienceBulk(ctx, id, studentIDs, lecturerID)
}

func (s *seminarService) UnvalidateAudience(ctx context.Context, id string, studentID string, lecturerID string) error {
	return s.seminarRepo.UnvalidateAudience(ctx, id, studentID, lecturerID)
}
