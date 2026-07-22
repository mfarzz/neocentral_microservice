package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"neocentral-go/notification-service/internal/model"
	"neocentral-go/notification-service/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo repository.NotificationRepository
	sse  *SSEManager
}

func NewNotificationService(repo repository.NotificationRepository, sse *SSEManager) *NotificationService {
	return &NotificationService{
		repo: repo,
		sse:  sse,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID, title, message string) error {
	notif := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Message:   &message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, notif)
	if err != nil {
		return err
	}

	// Push to SSE client if connected
	payload, err := json.Marshal(notif)
	if err == nil {
		s.sse.SendToUser(userID, payload)
	} else {
		log.Printf("Failed to marshal notification: %v", err)
	}

	return nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string) ([]model.Notification, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
