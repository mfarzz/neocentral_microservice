package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"gorm.io/gorm"
)

type ProposalRepository interface {
	Create(ctx context.Context, proposal *model.InternshipProposal) error
	GetByID(ctx context.Context, id string) (*model.InternshipProposal, error)
	FindByCoordinatorID(ctx context.Context, coordinatorID string) ([]model.InternshipProposal, error)
	Update(ctx context.Context, proposal *model.InternshipProposal) error
	Delete(ctx context.Context, id string) error
}

type proposalRepository struct {
	db *gorm.DB
}

func NewProposalRepository(db *gorm.DB) ProposalRepository {
	return &proposalRepository{db: db}
}

func (r *proposalRepository) Create(ctx context.Context, proposal *model.InternshipProposal) error {
	return r.db.WithContext(ctx).Create(proposal).Error
}

func (r *proposalRepository) GetByID(ctx context.Context, id string) (*model.InternshipProposal, error) {
	var proposal model.InternshipProposal
	if err := r.db.WithContext(ctx).Preload("Company").Where("id = ?", id).First(&proposal).Error; err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (r *proposalRepository) FindByCoordinatorID(ctx context.Context, coordinatorID string) ([]model.InternshipProposal, error) {
	var proposals []model.InternshipProposal
	if err := r.db.WithContext(ctx).Preload("Company").Where("coordinator_id = ?", coordinatorID).Find(&proposals).Error; err != nil {
		return nil, err
	}
	return proposals, nil
}

func (r *proposalRepository) Update(ctx context.Context, proposal *model.InternshipProposal) error {
	return r.db.WithContext(ctx).Save(proposal).Error
}

func (r *proposalRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.InternshipProposal{}).Error
}
