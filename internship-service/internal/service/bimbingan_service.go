package service

import (
	"context"

	"neocentral-go/internship-service/internal/model"
	"neocentral-go/internship-service/internal/repository"
)

type BimbinganService interface {
	GetQuestions(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceQuestion, error)
	CreateQuestion(ctx context.Context, question *model.InternshipGuidanceQuestion) error
	GetLecturerCriteria(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceLecturerCriteria, error)
	CreateLecturerCriteria(ctx context.Context, criteria *model.InternshipGuidanceLecturerCriteria) error
	
	GetStudentGuidanceSessions(ctx context.Context, studentID string) ([]model.InternshipGuidanceSession, error)
	SubmitGuidanceSession(ctx context.Context, session *model.InternshipGuidanceSession) error
	SubmitLecturerEvaluation(ctx context.Context, sessionID string, answers []model.InternshipGuidanceLecturerAnswer) error
}

type bimbinganService struct {
	bimbinganRepo  repository.BimbinganRepository
	internshipRepo repository.InternshipRepository
}

func NewBimbinganService(bimbinganRepo repository.BimbinganRepository, internshipRepo repository.InternshipRepository) BimbinganService {
	return &bimbinganService{
		bimbinganRepo:  bimbinganRepo,
		internshipRepo: internshipRepo,
	}
}

func (s *bimbinganService) GetQuestions(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceQuestion, error) {
	return s.bimbinganRepo.GetQuestionsByAcademicYear(ctx, academicYearID)
}

func (s *bimbinganService) CreateQuestion(ctx context.Context, question *model.InternshipGuidanceQuestion) error {
	return s.bimbinganRepo.CreateQuestion(ctx, question)
}

func (s *bimbinganService) GetLecturerCriteria(ctx context.Context, academicYearID string) ([]model.InternshipGuidanceLecturerCriteria, error) {
	return s.bimbinganRepo.GetLecturerCriteriaByAcademicYear(ctx, academicYearID)
}

func (s *bimbinganService) CreateLecturerCriteria(ctx context.Context, criteria *model.InternshipGuidanceLecturerCriteria) error {
	return s.bimbinganRepo.CreateLecturerCriteria(ctx, criteria)
}

func (s *bimbinganService) GetStudentGuidanceSessions(ctx context.Context, studentID string) ([]model.InternshipGuidanceSession, error) {
	internships, err := s.internshipRepo.FindByStudentID(ctx, studentID)
	if err != nil || len(internships) == 0 {
		return nil, err
	}
	internshipID := internships[0].ID
	return s.bimbinganRepo.GetSessionsByInternship(ctx, internshipID)
}

func (s *bimbinganService) SubmitGuidanceSession(ctx context.Context, session *model.InternshipGuidanceSession) error {
	return s.bimbinganRepo.CreateSessionWithAnswers(ctx, session)
}

func (s *bimbinganService) SubmitLecturerEvaluation(ctx context.Context, sessionID string, answers []model.InternshipGuidanceLecturerAnswer) error {
	return s.bimbinganRepo.SubmitLecturerEvaluation(ctx, sessionID, answers)
}
