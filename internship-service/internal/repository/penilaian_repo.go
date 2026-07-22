package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type PenilaianRepository interface {
	// CPMK
	GetCPMKsByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipCPMK, error)
	GetCPMKByID(ctx context.Context, id string) (*model.InternshipCPMK, error)
	CreateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error
	UpdateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error
	DeleteCPMK(ctx context.Context, id string) error

	// Rubrics
	CreateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error
	UpdateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error
	DeleteRubric(ctx context.Context, id string) error
	GetRubricsByCPMK(ctx context.Context, cpmkID string) ([]model.InternshipAssessmentRubric, error)
	GetRubricByID(ctx context.Context, id string) (*model.InternshipAssessmentRubric, error)

	// Assessments
	SaveLecturerScores(ctx context.Context, scores []model.InternshipLecturerScore) error
	GetLecturerScoresByInternship(ctx context.Context, internshipID string) ([]model.InternshipLecturerScore, error)
	SaveFieldScores(ctx context.Context, scores []model.InternshipFieldScore) error
	GetFieldScoresByInternship(ctx context.Context, internshipID string) ([]model.InternshipFieldScore, error)

	// Token
	GetFieldAssessmentToken(ctx context.Context, token string) (*model.FieldAssessmentToken, error)
	MarkTokenAsUsed(ctx context.Context, tokenID string) error
}

type penilaianRepository struct {
	db *gorm.DB
}

func NewPenilaianRepository(db *gorm.DB) PenilaianRepository {
	return &penilaianRepository{db: db}
}

func (r *penilaianRepository) GetCPMKsByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipCPMK, error) {
	var cpmks []model.InternshipCPMK
	err := r.db.WithContext(ctx).Preload("Rubrics").Where("academic_year_id = ?", academicYearID).Find(&cpmks).Error
	return cpmks, err
}

func (r *penilaianRepository) GetCPMKByID(ctx context.Context, id string) (*model.InternshipCPMK, error) {
	var cpmk model.InternshipCPMK
	if err := r.db.WithContext(ctx).Preload("Rubrics").Where("id = ?", id).First(&cpmk).Error; err != nil {
		return nil, err
	}
	return &cpmk, nil
}

func (r *penilaianRepository) CreateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error {
	return r.db.WithContext(ctx).Create(cpmk).Error
}

func (r *penilaianRepository) UpdateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error {
	return r.db.WithContext(ctx).Save(cpmk).Error
}

func (r *penilaianRepository) DeleteCPMK(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.InternshipCPMK{}).Error
}

func (r *penilaianRepository) CreateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error {
	return r.db.WithContext(ctx).Create(rubric).Error
}

func (r *penilaianRepository) UpdateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error {
	return r.db.WithContext(ctx).Save(rubric).Error
}

func (r *penilaianRepository) DeleteRubric(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.InternshipAssessmentRubric{}).Error
}

func (r *penilaianRepository) GetRubricsByCPMK(ctx context.Context, cpmkID string) ([]model.InternshipAssessmentRubric, error) {
	var rubrics []model.InternshipAssessmentRubric
	err := r.db.WithContext(ctx).Where("cpmk_id = ?", cpmkID).Find(&rubrics).Error
	return rubrics, err
}

func (r *penilaianRepository) GetRubricByID(ctx context.Context, id string) (*model.InternshipAssessmentRubric, error) {
	var rubric model.InternshipAssessmentRubric
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rubric).Error; err != nil {
		return nil, err
	}
	return &rubric, nil
}

func (r *penilaianRepository) SaveLecturerScores(ctx context.Context, scores []model.InternshipLecturerScore) error {
	if len(scores) == 0 {
		return nil
	}
	internshipID := scores[0].InternshipID

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("internship_id = ?", internshipID).Delete(&model.InternshipLecturerScore{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&scores).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *penilaianRepository) GetLecturerScoresByInternship(ctx context.Context, internshipID string) ([]model.InternshipLecturerScore, error) {
	var scores []model.InternshipLecturerScore
	err := r.db.WithContext(ctx).Where("internship_id = ?", internshipID).Find(&scores).Error
	return scores, err
}

func (r *penilaianRepository) SaveFieldScores(ctx context.Context, scores []model.InternshipFieldScore) error {
	if len(scores) == 0 {
		return nil
	}
	internshipID := scores[0].InternshipID

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("internship_id = ?", internshipID).Delete(&model.InternshipFieldScore{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&scores).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *penilaianRepository) GetFieldScoresByInternship(ctx context.Context, internshipID string) ([]model.InternshipFieldScore, error) {
	var scores []model.InternshipFieldScore
	err := r.db.WithContext(ctx).Where("internship_id = ?", internshipID).Find(&scores).Error
	return scores, err
}

func (r *penilaianRepository) GetFieldAssessmentToken(ctx context.Context, token string) (*model.FieldAssessmentToken, error) {
	var fieldAssessmentToken model.FieldAssessmentToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&fieldAssessmentToken).Error; err != nil {
		return nil, err
	}
	return &fieldAssessmentToken, nil
}

func (r *penilaianRepository) MarkTokenAsUsed(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).Model(&model.FieldAssessmentToken{}).Where("id = ?", tokenID).Update("is_used", true).Error
}
