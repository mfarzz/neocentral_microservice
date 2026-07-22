package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type BimbinganRepository interface {
	// Questions
	GetQuestionsByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceQuestion, error)
	CreateQuestion(ctx context.Context, question *model.InternshipGuidanceQuestion) error
	
	// Lecturer Criteria
	GetLecturerCriteriaByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceLecturerCriteria, error)
	CreateLecturerCriteria(ctx context.Context, criteria *model.InternshipGuidanceLecturerCriteria) error

	// Sessions & Answers
	GetSessionsByInternship(ctx context.Context, internshipID string) ([]model.InternshipGuidanceSession, error)
	CreateSessionWithAnswers(ctx context.Context, session *model.InternshipGuidanceSession) error
	SubmitLecturerEvaluation(ctx context.Context, sessionID string, answers []model.InternshipGuidanceLecturerAnswer) error
}

type bimbinganRepository struct {
	db *gorm.DB
}

func NewBimbinganRepository(db *gorm.DB) BimbinganRepository {
	return &bimbinganRepository{db: db}
}

func (r *bimbinganRepository) GetQuestionsByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceQuestion, error) {
	var questions []model.InternshipGuidanceQuestion
	err := r.db.WithContext(ctx).Where("academic_year_id = ?", academicYearID).Find(&questions).Error
	return questions, err
}

func (r *bimbinganRepository) CreateQuestion(ctx context.Context, question *model.InternshipGuidanceQuestion) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *bimbinganRepository) GetLecturerCriteriaByAcademicYear(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceLecturerCriteria, error) {
	var criteria []model.InternshipGuidanceLecturerCriteria
	err := r.db.WithContext(ctx).Preload("Options").Where("academic_year_id = ?", academicYearID).Find(&criteria).Error
	return criteria, err
}

func (r *bimbinganRepository) CreateLecturerCriteria(ctx context.Context, criteria *model.InternshipGuidanceLecturerCriteria) error {
	return r.db.WithContext(ctx).Create(criteria).Error
}

func (r *bimbinganRepository) GetSessionsByInternship(ctx context.Context, internshipID string) ([]model.InternshipGuidanceSession, error) {
	var sessions []model.InternshipGuidanceSession
	err := r.db.WithContext(ctx).
		Preload("StudentAnswers").
		Preload("LecturerAnswers").
		Where("internship_id = ?", internshipID).
		Order("week_number ASC").Find(&sessions).Error
	return sessions, err
}

func (r *bimbinganRepository) CreateSessionWithAnswers(ctx context.Context, session *model.InternshipGuidanceSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(session).Error; err != nil {
			return err
		}
		// Assuming StudentAnswers are attached to session and GORM creates them automatically
		return nil
	})
}

func (r *bimbinganRepository) SubmitLecturerEvaluation(ctx context.Context, sessionID string, answers []model.InternshipGuidanceLecturerAnswer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&answers).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.InternshipGuidanceSession{}).
			Where("id = ?", sessionID).
			Update("status", model.SessionApproved).Error; err != nil {
			return err
		}
		return nil
	})
}
