package repository

import (
	"context"

	"gorm.io/gorm"
	"neocentral-go/master-data-service/internal/model"
)

type gormRoomRepo struct{ db *gorm.DB }

func NewGormRoomRepo(db *gorm.DB) RoomRepository {
	return &gormRoomRepo{db: db}
}

func (r *gormRoomRepo) FindAll(ctx context.Context) ([]model.Room, error) {
	var list []model.Room
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormRoomRepo) FindByID(ctx context.Context, id string) (*model.Room, error) {
	var room model.Room
	if err := r.db.WithContext(ctx).First(&room, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *gormRoomRepo) Create(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *gormRoomRepo) Update(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *gormRoomRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Room{}, "id = ?", id).Error
}
