package service

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
	"github.com/google/uuid"
)

type PendaftaranService interface {
	CreateProposal(ctx context.Context, proposal *model.InternshipProposal) error
	GetProposals(ctx context.Context) ([]model.InternshipProposal, error)
	GetProposalByID(ctx context.Context, id string) (*model.InternshipProposal, error)
	ListStudentProposals(ctx context.Context, studentID string) ([]model.InternshipProposal, error)
	UpdateProposal(ctx context.Context, id string, updatedData *model.InternshipProposal) error
	UpdateProposalStatus(ctx context.Context, id string, status model.ProposalStatus) error
	DeleteProposal(ctx context.Context, id string) error
	ListCompanies(ctx context.Context) ([]model.Company, error)
}

type pendaftaranService struct {
	proposalRepo repository.ProposalRepository
	companyRepo  repository.CompanyRepository
}

func NewPendaftaranService(proposalRepo repository.ProposalRepository, companyRepo repository.CompanyRepository) PendaftaranService {
	return &pendaftaranService{
		proposalRepo: proposalRepo,
		companyRepo:  companyRepo,
	}
}

func (s *pendaftaranService) CreateProposal(ctx context.Context, proposal *model.InternshipProposal) error {
	proposal.ID = uuid.New().String()
	proposal.Status = model.ProposalPending
	return s.proposalRepo.Create(ctx, proposal)
}

func (s *pendaftaranService) GetProposals(ctx context.Context) ([]model.InternshipProposal, error) {
	// Usually this gets all proposals for an admin
	// Since there's no ListAll in repo yet, we'll return an empty list or implement ListAll
	return []model.InternshipProposal{}, nil
}

func (s *pendaftaranService) GetProposalByID(ctx context.Context, id string) (*model.InternshipProposal, error) {
	return s.proposalRepo.GetByID(ctx, id)
}

func (s *pendaftaranService) ListStudentProposals(ctx context.Context, studentID string) ([]model.InternshipProposal, error) {
	return s.proposalRepo.FindByCoordinatorID(ctx, studentID)
}

func (s *pendaftaranService) UpdateProposal(ctx context.Context, id string, updatedData *model.InternshipProposal) error {
	existing, err := s.proposalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.TargetCompanyID = updatedData.TargetCompanyID
	existing.ProposalDocumentID = updatedData.ProposalDocumentID
	existing.ProposedStartDate = updatedData.ProposedStartDate
	existing.ProposedEndDate = updatedData.ProposedEndDate

	return s.proposalRepo.Update(ctx, existing)
}

func (s *pendaftaranService) UpdateProposalStatus(ctx context.Context, id string, status model.ProposalStatus) error {
	existing, err := s.proposalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.Status = status
	return s.proposalRepo.Update(ctx, existing)
}

func (s *pendaftaranService) DeleteProposal(ctx context.Context, id string) error {
	// Optional: Check if proposal exists and user has permission
	_, err := s.proposalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.proposalRepo.Delete(ctx, id)
}

func (s *pendaftaranService) ListCompanies(ctx context.Context) ([]model.Company, error) {
	return s.companyRepo.FindAll(ctx)
}
