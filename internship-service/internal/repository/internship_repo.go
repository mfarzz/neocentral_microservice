package repository

import (
	"context"

	"neocentral-go/internship-service/internal/model"

	"gorm.io/gorm"
)

type InternshipRepository interface {
	Create(ctx context.Context, internship *model.Internship) error
	GetByID(ctx context.Context, id string) (*model.Internship, error)
	FindByStudentID(ctx context.Context, studentID string) ([]model.Internship, error)
	Update(ctx context.Context, internship *model.Internship) error
	UpdateSupervisorBulk(ctx context.Context, internshipIDs []string, supervisorID string) error
	GetLecturerWorkload(ctx context.Context, academicYearID string) ([]map[string]interface{}, error)
	List(ctx context.Context) ([]model.Internship, int64, error)
}

type internshipRepository struct {
	db *gorm.DB
}

func NewInternshipRepository(db *gorm.DB) InternshipRepository {
	return &internshipRepository{db: db}
}

func (r *internshipRepository) Create(ctx context.Context, internship *model.Internship) error {
	return r.db.WithContext(ctx).Create(internship).Error
}

func (r *internshipRepository) GetByID(ctx context.Context, id string) (*model.Internship, error) {
	var internship model.Internship
	if err := r.db.WithContext(ctx).Preload("Proposal").Preload("Proposal.Company").Where("id = ?", id).First(&internship).Error; err != nil {
		return nil, err
	}
	return &internship, nil
}

func (r *internshipRepository) FindByStudentID(ctx context.Context, studentID string) ([]model.Internship, error) {
	var internships []model.Internship
	if err := r.db.WithContext(ctx).Preload("Proposal").Preload("Proposal.Company").Where("student_id = ?", studentID).Find(&internships).Error; err != nil {
		return nil, err
	}
	return internships, nil
}

func (r *internshipRepository) Update(ctx context.Context, internship *model.Internship) error {
	return r.db.WithContext(ctx).Save(internship).Error
}

func (r *internshipRepository) UpdateSupervisorBulk(ctx context.Context, internshipIDs []string, supervisorID string) error {
	return r.db.WithContext(ctx).Model(&model.Internship{}).Where("id IN ?", internshipIDs).Update("supervisor_id", supervisorID).Error
}

func (r *internshipRepository) GetLecturerWorkload(ctx context.Context, academicYearID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	query := r.db.WithContext(ctx).
		Table("internships").
		Select("supervisor_id, COUNT(id) as total_internships").
		Where("supervisor_id IS NOT NULL")
		
	// Optional filter by academic year through the proposal
	if academicYearID != "" {
		query = query.Joins("JOIN internship_proposals ip ON internships.proposal_id = ip.id").
			Where("ip.academic_year_id = ?", academicYearID)
	}
	
	err := query.Group("supervisor_id").Find(&results).Error
	return results, err
}

func (r *internshipRepository) List(ctx context.Context) ([]model.Internship, int64, error) {
	var internships []model.Internship
	var count int64
	
	err := r.db.WithContext(ctx).Model(&model.Internship{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = r.db.WithContext(ctx).Preload("Proposal").Preload("Proposal.Company").Find(&internships).Error
	return internships, count, err
}
