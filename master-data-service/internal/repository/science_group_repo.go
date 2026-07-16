package repository

import (
	"context"

	"gorm.io/gorm"
	"neocentral-go/master-data-service/internal/model"
)

type gormScienceGroupRepo struct{ db *gorm.DB }

func NewGormScienceGroupRepo(db *gorm.DB) ScienceGroupRepository {
	return &gormScienceGroupRepo{db: db}
}

func (r *gormScienceGroupRepo) FindAll(ctx context.Context) ([]model.ScienceGroup, error) {
	var list []model.ScienceGroup
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormScienceGroupRepo) FindByID(ctx context.Context, id string) (*model.ScienceGroup, error) {
	var sg model.ScienceGroup
	if err := r.db.WithContext(ctx).First(&sg, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sg, nil
}

func (r *gormScienceGroupRepo) Create(ctx context.Context, sg *model.ScienceGroup) error {
	return r.db.WithContext(ctx).Create(sg).Error
}

func (r *gormScienceGroupRepo) Update(ctx context.Context, sg *model.ScienceGroup) error {
	return r.db.WithContext(ctx).Save(sg).Error
}

func (r *gormScienceGroupRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.ScienceGroup{}, "id = ?", id).Error
}
