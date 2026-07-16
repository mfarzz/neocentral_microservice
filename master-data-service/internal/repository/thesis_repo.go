package repository

import (
	"context"

	"gorm.io/gorm"
	"neocentral-go/master-data-service/internal/model"
)

type gormThesisTopicRepo struct{ db *gorm.DB }

func NewGormThesisTopicRepo(db *gorm.DB) ThesisTopicRepository {
	return &gormThesisTopicRepo{db: db}
}

func (r *gormThesisTopicRepo) FindAll(ctx context.Context) ([]model.ThesisTopic, error) {
	var list []model.ThesisTopic
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormThesisTopicRepo) FindByID(ctx context.Context, id string) (*model.ThesisTopic, error) {
	var topic model.ThesisTopic
	if err := r.db.WithContext(ctx).First(&topic, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

func (r *gormThesisTopicRepo) Create(ctx context.Context, topic *model.ThesisTopic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *gormThesisTopicRepo) Update(ctx context.Context, topic *model.ThesisTopic) error {
	return r.db.WithContext(ctx).Save(topic).Error
}

func (r *gormThesisTopicRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.ThesisTopic{}, "id = ?", id).Error
}

// ── Thesis Status Repository ────────────────

type gormThesisStatusRepo struct{ db *gorm.DB }

func NewGormThesisStatusRepo(db *gorm.DB) ThesisStatusRepository {
	return &gormThesisStatusRepo{db: db}
}

func (r *gormThesisStatusRepo) FindAll(ctx context.Context) ([]model.ThesisStatus, error) {
	var list []model.ThesisStatus
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormThesisStatusRepo) FindByID(ctx context.Context, id string) (*model.ThesisStatus, error) {
	var status model.ThesisStatus
	if err := r.db.WithContext(ctx).First(&status, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}
