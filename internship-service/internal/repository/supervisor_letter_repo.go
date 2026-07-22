package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type SupervisorLetterRepository interface {
	GetBySupervisorID(ctx context.Context, supervisorID string) (*model.InternshipSupervisorLetter, error)
	Create(ctx context.Context, letter *model.InternshipSupervisorLetter) error
	Update(ctx context.Context, letter *model.InternshipSupervisorLetter) error
}

type supervisorLetterRepository struct {
	db *gorm.DB
}

func NewSupervisorLetterRepository(db *gorm.DB) SupervisorLetterRepository {
	return &supervisorLetterRepository{db: db}
}

func (r *supervisorLetterRepository) GetBySupervisorID(ctx context.Context, supervisorID string) (*model.InternshipSupervisorLetter, error) {
	var letter model.InternshipSupervisorLetter
	if err := r.db.WithContext(ctx).Where("supervisor_id = ?", supervisorID).First(&letter).Error; err != nil {
		return nil, err
	}
	return &letter, nil
}

func (r *supervisorLetterRepository) Create(ctx context.Context, letter *model.InternshipSupervisorLetter) error {
	return r.db.WithContext(ctx).Create(letter).Error
}

func (r *supervisorLetterRepository) Update(ctx context.Context, letter *model.InternshipSupervisorLetter) error {
	return r.db.WithContext(ctx).Save(letter).Error
}
