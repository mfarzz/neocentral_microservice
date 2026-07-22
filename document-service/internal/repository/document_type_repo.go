package repository

import (
	"context"
	"neocentral-go/document-service/internal/model"

	"gorm.io/gorm"
)

type DocumentTypeRepository interface {
	Create(ctx context.Context, docType *model.DocumentType) error
	FindAll(ctx context.Context) ([]model.DocumentType, error)
}

type gormDocumentTypeRepo struct {
	db *gorm.DB
}

func NewDocumentTypeRepository(db *gorm.DB) DocumentTypeRepository {
	return &gormDocumentTypeRepo{db: db}
}

func (r *gormDocumentTypeRepo) Create(ctx context.Context, docType *model.DocumentType) error {
	return r.db.WithContext(ctx).Create(docType).Error
}

func (r *gormDocumentTypeRepo) FindAll(ctx context.Context) ([]model.DocumentType, error) {
	var types []model.DocumentType
	err := r.db.WithContext(ctx).Find(&types).Error
	return types, err
}
