package repository

import (
	"context"

	"neocentral-go/notification-service/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	FindByUserID(ctx context.Context, userID string) ([]model.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type gormNotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &gormNotificationRepo{db: db}
}

func (r *gormNotificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *gormNotificationRepo) FindByUserID(ctx context.Context, userID string) ([]model.Notification, error) {
	var notifs []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifs).Error
	return notifs, err
}

func (r *gormNotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *gormNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}
