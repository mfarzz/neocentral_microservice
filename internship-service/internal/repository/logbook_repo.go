package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type LogbookRepository interface {
	GetByInternshipID(ctx context.Context, internshipID string) ([]model.InternshipLogbook, error)
	Create(ctx context.Context, logbook *model.InternshipLogbook) error
	Update(ctx context.Context, logbook *model.InternshipLogbook) error
	GetByID(ctx context.Context, id string) (*model.InternshipLogbook, error)
}

type logbookRepository struct {
	db *gorm.DB
}

func NewLogbookRepository(db *gorm.DB) LogbookRepository {
	return &logbookRepository{db: db}
}

func (r *logbookRepository) GetByInternshipID(ctx context.Context, internshipID string) ([]model.InternshipLogbook, error) {
	var logbooks []model.InternshipLogbook
	if err := r.db.WithContext(ctx).Where("internship_id = ?", internshipID).Order("activity_date ASC").Find(&logbooks).Error; err != nil {
		return nil, err
	}
	return logbooks, nil
}

func (r *logbookRepository) Create(ctx context.Context, logbook *model.InternshipLogbook) error {
	return r.db.WithContext(ctx).Create(logbook).Error
}

func (r *logbookRepository) Update(ctx context.Context, logbook *model.InternshipLogbook) error {
	return r.db.WithContext(ctx).Save(logbook).Error
}

func (r *logbookRepository) GetByID(ctx context.Context, id string) (*model.InternshipLogbook, error) {
	var logbook model.InternshipLogbook
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&logbook).Error; err != nil {
		return nil, err
	}
	return &logbook, nil
}
