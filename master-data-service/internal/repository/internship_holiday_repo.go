package repository

import (
	"context"

	"gorm.io/gorm"
	"neocentral-go/master-data-service/internal/model"
)

type InternshipHolidayRepository interface {
	FindAll(ctx context.Context) ([]model.InternshipHoliday, error)
	FindByID(ctx context.Context, id string) (*model.InternshipHoliday, error)
	FindByDate(ctx context.Context, date string) (*model.InternshipHoliday, error)
	Create(ctx context.Context, holiday *model.InternshipHoliday) error
	Update(ctx context.Context, holiday *model.InternshipHoliday) error
	Delete(ctx context.Context, id string) error
}

type gormInternshipHolidayRepo struct {
	db *gorm.DB
}

func NewGormInternshipHolidayRepo(db *gorm.DB) InternshipHolidayRepository {
	return &gormInternshipHolidayRepo{db: db}
}

func (r *gormInternshipHolidayRepo) FindAll(ctx context.Context) ([]model.InternshipHoliday, error) {
	var holidays []model.InternshipHoliday
	err := r.db.WithContext(ctx).Order("holiday_date desc").Find(&holidays).Error
	return holidays, err
}

func (r *gormInternshipHolidayRepo) FindByID(ctx context.Context, id string) (*model.InternshipHoliday, error) {
	var holiday model.InternshipHoliday
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&holiday).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &holiday, nil
}

func (r *gormInternshipHolidayRepo) FindByDate(ctx context.Context, date string) (*model.InternshipHoliday, error) {
	var holiday model.InternshipHoliday
	err := r.db.WithContext(ctx).Where("holiday_date = ?", date).First(&holiday).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &holiday, nil
}

func (r *gormInternshipHolidayRepo) Create(ctx context.Context, holiday *model.InternshipHoliday) error {
	return r.db.WithContext(ctx).Create(holiday).Error
}

func (r *gormInternshipHolidayRepo) Update(ctx context.Context, holiday *model.InternshipHoliday) error {
	return r.db.WithContext(ctx).Save(holiday).Error
}

func (r *gormInternshipHolidayRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.InternshipHoliday{}, "id = ?", id).Error
}
