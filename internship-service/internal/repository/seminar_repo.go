package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type SeminarRepository interface {
	GetUpcoming(ctx context.Context) ([]model.InternshipSeminar, error)
	GetByID(ctx context.Context, id string) (*model.InternshipSeminar, error)
	Create(ctx context.Context, seminar *model.InternshipSeminar) error
	Update(ctx context.Context, seminar *model.InternshipSeminar) error
	
	// Audience
	AddAudience(ctx context.Context, audience *model.InternshipSeminarAudience) error
	RemoveAudience(ctx context.Context, seminarID string, studentID string) error
	ValidateAudienceBulk(ctx context.Context, seminarID string, studentIDs []string, validatedBy string) error
	UnvalidateAudience(ctx context.Context, seminarID string, studentID string, validatedBy string) error
}

type seminarRepository struct {
	db *gorm.DB
}

func NewSeminarRepository(db *gorm.DB) SeminarRepository {
	return &seminarRepository{db: db}
}

func (r *seminarRepository) GetUpcoming(ctx context.Context) ([]model.InternshipSeminar, error) {
	var seminars []model.InternshipSeminar
	err := r.db.WithContext(ctx).Where("status = ?", model.SeminarApproved).Find(&seminars).Error
	return seminars, err
}

func (r *seminarRepository) GetByID(ctx context.Context, id string) (*model.InternshipSeminar, error) {
	var seminar model.InternshipSeminar
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&seminar).Error; err != nil {
		return nil, err
	}
	return &seminar, nil
}

func (r *seminarRepository) Create(ctx context.Context, seminar *model.InternshipSeminar) error {
	return r.db.WithContext(ctx).Create(seminar).Error
}

func (r *seminarRepository) Update(ctx context.Context, seminar *model.InternshipSeminar) error {
	return r.db.WithContext(ctx).Save(seminar).Error
}

func (r *seminarRepository) AddAudience(ctx context.Context, audience *model.InternshipSeminarAudience) error {
	return r.db.WithContext(ctx).Create(audience).Error
}

func (r *seminarRepository) RemoveAudience(ctx context.Context, seminarID string, studentID string) error {
	return r.db.WithContext(ctx).Where("seminar_id = ? AND student_id = ?", seminarID, studentID).Delete(&model.InternshipSeminarAudience{}).Error
}

func (r *seminarRepository) ValidateAudienceBulk(ctx context.Context, seminarID string, studentIDs []string, validatedBy string) error {
	return r.db.WithContext(ctx).Model(&model.InternshipSeminarAudience{}).
		Where("seminar_id = ? AND student_id IN ?", seminarID, studentIDs).
		Updates(map[string]interface{}{
			"is_validated": true,
			"validated_by": validatedBy,
		}).Error
}

func (r *seminarRepository) UnvalidateAudience(ctx context.Context, seminarID string, studentID string, validatedBy string) error {
	return r.db.WithContext(ctx).Model(&model.InternshipSeminarAudience{}).
		Where("seminar_id = ? AND student_id = ?", seminarID, studentID).
		Updates(map[string]interface{}{
			"is_validated": false,
			"validated_by": validatedBy, // keeping track of who unvalidated it
		}).Error
}
