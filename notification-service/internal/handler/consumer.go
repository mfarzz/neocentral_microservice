package handler

import (
	"context"
	"encoding/json"
	"log"

	"neocentral-go/notification-service/internal/service"
	"neocentral-go/notification-service/pkg/broker"

	"github.com/nats-io/nats.go"
)

type NotificationEvent struct {
	UserID  string `json:"userId"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type NatsConsumer struct {
	svc *service.NotificationService
}

func NewNatsConsumer(svc *service.NotificationService) *NatsConsumer {
	return &NatsConsumer{svc: svc}
}

func (c *NatsConsumer) Start(b *broker.NATSBroker) {
	_, err := b.Subscribe("notification.create", func(msg *nats.Msg) {
		var event NotificationEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal notification event: %v", err)
			return
		}

		err := c.svc.CreateNotification(context.Background(), event.UserID, event.Title, event.Message)
		if err != nil {
			log.Printf("Failed to create notification from event: %v", err)
		} else {
			log.Printf("Successfully processed notification event for user %s", event.UserID)
		}
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to notification.create: %v", err)
	}
}
