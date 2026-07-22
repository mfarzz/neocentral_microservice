package service

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type PelaksanaanService interface {
	GetInternshipHistory(ctx context.Context, studentID string) ([]model.Internship, error)
	
	// Logbook
	GetLogbooks(ctx context.Context, internshipID string) ([]model.InternshipLogbook, error)
	CreateLogbook(ctx context.Context, logbook *model.InternshipLogbook) error
	UpdateLogbook(ctx context.Context, id string, description string) error
	
	// Internship Updates
	UpdateInternshipDetails(ctx context.Context, internshipID string, data map[string]interface{}) error
	SubmitReport(ctx context.Context, internshipID string, title string, documentID string) error
	UpdateCompletionCertificate(ctx context.Context, internshipID string, documentID string) error
	UpdateCompanyReceipt(ctx context.Context, internshipID string, documentID string) error
	SubmitCompanyReport(ctx context.Context, internshipID string, documentID string) error
	SubmitLogbookDocument(ctx context.Context, internshipID string, documentID string) error
}

type pelaksanaanService struct {
	internshipRepo repository.InternshipRepository
	logbookRepo    repository.LogbookRepository
}

func NewPelaksanaanService(internshipRepo repository.InternshipRepository, logbookRepo repository.LogbookRepository) PelaksanaanService {
	return &pelaksanaanService{
		internshipRepo: internshipRepo,
		logbookRepo:    logbookRepo,
	}
}

func (s *pelaksanaanService) GetInternshipHistory(ctx context.Context, studentID string) ([]model.Internship, error) {
	return s.internshipRepo.FindByStudentID(ctx, studentID)
}

func (s *pelaksanaanService) GetLogbooks(ctx context.Context, internshipID string) ([]model.InternshipLogbook, error) {
	return s.logbookRepo.GetByInternshipID(ctx, internshipID)
}

func (s *pelaksanaanService) CreateLogbook(ctx context.Context, logbook *model.InternshipLogbook) error {
	return s.logbookRepo.Create(ctx, logbook)
}

func (s *pelaksanaanService) UpdateLogbook(ctx context.Context, id string, description string) error {
	logbook, err := s.logbookRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	logbook.ActivityDescription = description
	return s.logbookRepo.Update(ctx, logbook)
}

func (s *pelaksanaanService) UpdateInternshipDetails(ctx context.Context, internshipID string, data map[string]interface{}) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	if name, ok := data["fieldSupervisorName"].(string); ok {
		internship.FieldSupervisorName = &name
	}
	if email, ok := data["fieldSupervisorEmail"].(string); ok {
		internship.FieldSupervisorEmail = &email
	}
	if phone, ok := data["fieldSupervisorPhone"].(string); ok {
		internship.FieldSupervisorPhone = &phone
	}
	if nip, ok := data["fieldSupervisorNip"].(string); ok {
		internship.FieldSupervisorNip = &nip
	}
	if unit, ok := data["unitSection"].(string); ok {
		internship.UnitSection = &unit
	}
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *pelaksanaanService) SubmitReport(ctx context.Context, internshipID string, title string, documentID string) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	status := model.DocSubmitted
	internship.ReportTitle = &title
	internship.ReportDocumentID = &documentID
	internship.ReportStatus = &status
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *pelaksanaanService) UpdateCompletionCertificate(ctx context.Context, internshipID string, documentID string) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	status := model.DocSubmitted
	internship.CompletionCertificateDocID = &documentID
	internship.CompletionCertificateStatus = &status
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *pelaksanaanService) UpdateCompanyReceipt(ctx context.Context, internshipID string, documentID string) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	status := model.DocSubmitted
	internship.CompanyReceiptDocID = &documentID
	internship.CompanyReceiptStatus = &status
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *pelaksanaanService) SubmitCompanyReport(ctx context.Context, internshipID string, documentID string) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	status := model.DocSubmitted
	internship.CompanyReportDocID = &documentID
	internship.CompanyReportStatus = &status
	
	return s.internshipRepo.Update(ctx, internship)
}

func (s *pelaksanaanService) SubmitLogbookDocument(ctx context.Context, internshipID string, documentID string) error {
	internship, err := s.internshipRepo.GetByID(ctx, internshipID)
	if err != nil {
		return err
	}
	
	status := model.DocSubmitted
	internship.LogbookDocumentID = &documentID
	internship.LogbookDocumentStatus = &status
	
	return s.internshipRepo.Update(ctx, internship)
}
