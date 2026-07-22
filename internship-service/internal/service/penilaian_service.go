package service

import (
	"context"
	"errors"
	"time"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type ScoreInput struct {
	RubricID string  `json:"rubricId"`
	Score    float64 `json:"score"`
}

type PenilaianService interface {
	// CPMK & Rubrics
	GetCPMKs(ctx context.Context, academicYearID string) ([]model.InternshipCPMK, error)
	CreateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error
	CreateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error
	
	// Assessments
	SubmitLecturerAssessment(ctx context.Context, internshipID string, scores []ScoreInput) error
	SubmitFieldAssessment(ctx context.Context, token string, scores []ScoreInput, notes string, signature string) error
}

type penilaianService struct {
	penilaianRepo repository.PenilaianRepository
}

func NewPenilaianService(penilaianRepo repository.PenilaianRepository) PenilaianService {
	return &penilaianService{
		penilaianRepo: penilaianRepo,
	}
}

func (s *penilaianService) GetCPMKs(ctx context.Context, academicYearID string) ([]model.InternshipCPMK, error) {
	return s.penilaianRepo.GetCPMKsByAcademicYear(ctx, academicYearID)
}

func (s *penilaianService) CreateCPMK(ctx context.Context, cpmk *model.InternshipCPMK) error {
	return s.penilaianRepo.CreateCPMK(ctx, cpmk)
}

func (s *penilaianService) CreateRubric(ctx context.Context, rubric *model.InternshipAssessmentRubric) error {
	return s.penilaianRepo.CreateRubric(ctx, rubric)
}

func (s *penilaianService) SubmitLecturerAssessment(ctx context.Context, internshipID string, scores []ScoreInput) error {
	var lecturerScores []model.InternshipLecturerScore
	for _, sc := range scores {
		lecturerScores = append(lecturerScores, model.InternshipLecturerScore{
			InternshipID:   internshipID,
			ChosenRubricID: sc.RubricID,
			Score:          sc.Score,
		})
	}
	
	return s.penilaianRepo.SaveLecturerScores(ctx, lecturerScores)
}

func (s *penilaianService) SubmitFieldAssessment(ctx context.Context, token string, scores []ScoreInput, notes string, signature string) error {
	// 1. Validate token
	fieldToken, err := s.penilaianRepo.GetFieldAssessmentToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired field assessment token")
	}

	if fieldToken.IsUsed {
		return errors.New("field assessment token has already been used")
	}

	if fieldToken.ExpiresAt.Before(time.Now()) {
		return errors.New("field assessment token has expired")
	}

	internshipID := fieldToken.InternshipID

	var fieldScores []model.InternshipFieldScore
	for _, sc := range scores {
		fieldScores = append(fieldScores, model.InternshipFieldScore{
			InternshipID:   internshipID,
			ChosenRubricID: sc.RubricID,
			Score:          sc.Score,
		})
	}

	// Save assessments
	if err := s.penilaianRepo.SaveFieldScores(ctx, fieldScores); err != nil {
		return err
	}
	
	// Mark token as used
	if err := s.penilaianRepo.MarkTokenAsUsed(ctx, fieldToken.ID); err != nil {
		return err
	}

	// Notes and signature are typically saved in the Internship table (FieldAssessmentNotes, etc.)
	// You might want to trigger an update to the Internship record here.

	return nil
}
