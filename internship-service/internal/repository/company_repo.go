package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	FindAll(ctx context.Context) ([]model.Company, error)
	GetByID(ctx context.Context, id string) (*model.Company, error)
	Create(ctx context.Context, company *model.Company) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) FindAll(ctx context.Context) ([]model.Company, error) {
	var companies []model.Company
	if err := r.db.WithContext(ctx).Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) GetByID(ctx context.Context, id string) (*model.Company, error) {
	var company model.Company
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Create(ctx context.Context, company *model.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}
