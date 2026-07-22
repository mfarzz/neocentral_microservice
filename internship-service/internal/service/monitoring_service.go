package service

import (
	"context"
	
	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type MonitoringService interface {
	GetInternshipList(ctx context.Context, filters map[string]interface{}, page int, pageSize int) ([]model.Internship, int64, error)
	GetInternshipDetail(ctx context.Context, id string) (*model.Internship, error)
	GetGradeRecap(ctx context.Context, academicYearID string) (interface{}, error)
	VerifyDocument(ctx context.Context, id string, docType string, status string, notes string) error
	RejectFinalReport(ctx context.Context, id string, notes string) error
	GetMonitoringStats(ctx context.Context, academicYearID string) (map[string]interface{}, error)
}

type monitoringService struct {
	internshipRepo repository.InternshipRepository
}

func NewMonitoringService(internshipRepo repository.InternshipRepository) MonitoringService {
	return &monitoringService{
		internshipRepo: internshipRepo,
	}
}

func (s *monitoringService) GetInternshipList(ctx context.Context, filters map[string]interface{}, page int, pageSize int) ([]model.Internship, int64, error) {
	// Simplify: just fetch all from repo and count (in real life we use a specific repo method with filters/pagination)
	return s.internshipRepo.List(ctx)
}

func (s *monitoringService) GetInternshipDetail(ctx context.Context, id string) (*model.Internship, error) {
	return s.internshipRepo.GetByID(ctx, id)
}

func (s *monitoringService) GetGradeRecap(ctx context.Context, academicYearID string) (interface{}, error) {
	// Usually this aggregates CPMK scores and internships based on the academic year.
	// We return a dummy implementation for the skeleton.
	return map[string]interface{}{
		"cpmks": []string{},
		"items": []string{},
	}, nil
}

func (s *monitoringService) VerifyDocument(ctx context.Context, id string, docType string, status string, notes string) error {
	internship, err := s.internshipRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Example for completion certificate
	if docType == "completion_certificate" {
		docStatus := model.DocumentSubmitStatus(status)
		internship.CompletionCertificateStatus = &docStatus
		internship.CompletionCertificateNotes = &notes
	}
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *monitoringService) RejectFinalReport(ctx context.Context, id string, notes string) error {
	internship, err := s.internshipRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	rejectedStatus := model.DocumentSubmitStatus("rejected")
	internship.ReportStatus = &rejectedStatus
	internship.ReportNotes = &notes
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *monitoringService) GetMonitoringStats(ctx context.Context, academicYearID string) (map[string]interface{}, error) {
	// Count ongoing, completed, etc.
	// Returning dummy for skeleton.
	return map[string]interface{}{
		"total": 0,
		"ongoing": 0,
		"completed": 0,
	}, nil
}
