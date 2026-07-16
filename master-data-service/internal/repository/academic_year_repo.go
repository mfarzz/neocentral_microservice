package repository

import (
	"context"

	"gorm.io/gorm"
	"neocentral-go/master-data-service/internal/model"
)

type gormAcademicYearRepo struct{ db *gorm.DB }

func NewGormAcademicYearRepo(db *gorm.DB) AcademicYearRepository {
	return &gormAcademicYearRepo{db: db}
}

func (r *gormAcademicYearRepo) FindAll(ctx context.Context) ([]model.AcademicYear, error) {
	var list []model.AcademicYear
	if err := r.db.WithContext(ctx).Order("year DESC, semester ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormAcademicYearRepo) FindByID(ctx context.Context, id string) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	if err := r.db.WithContext(ctx).First(&ay, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ay, nil
}

func (r *gormAcademicYearRepo) FindActive(ctx context.Context) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&ay).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ay, nil
}

func (r *gormAcademicYearRepo) Create(ctx context.Context, ay *model.AcademicYear) error {
	return r.db.WithContext(ctx).Create(ay).Error
}

func (r *gormAcademicYearRepo) Update(ctx context.Context, ay *model.AcademicYear) error {
	return r.db.WithContext(ctx).Save(ay).Error
}

func (r *gormAcademicYearRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.AcademicYear{}, "id = ?", id).Error
}

func (r *gormAcademicYearRepo) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&model.AcademicYear{}).Where("is_active = ?", true).Update("is_active", false).Error
}
