package repository

import (
	"context"

	"neocentral-go/document-service/internal/model"

	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *model.Document) error
	FindByID(ctx context.Context, id string) (*model.Document, error)
	Delete(ctx context.Context, id string) error
}

type gormDocumentRepo struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &gormDocumentRepo{db: db}
}

func (r *gormDocumentRepo) Create(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *gormDocumentRepo) FindByID(ctx context.Context, id string) (*model.Document, error) {
	var doc model.Document
	err := r.db.WithContext(ctx).First(&doc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *gormDocumentRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Document{}, "id = ?", id).Error
}
