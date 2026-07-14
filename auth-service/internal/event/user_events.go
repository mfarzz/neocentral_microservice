package event

import (
	"time"

	"neocentral-go/pkg/messaging"
)

// UserEventPublisher publishes user-related events to NATS.
type UserEventPublisher struct {
	nats *messaging.NATSClient
}

func NewUserEventPublisher(nats *messaging.NATSClient) *UserEventPublisher {
	return &UserEventPublisher{nats: nats}
}

// ── Event payloads ───────────────────────────────────────────────

type UserCreatedEvent struct {
	UserID         string `json:"userId"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
	IdentityType   string `json:"identityType"`
	IdentityNumber string `json:"identityNumber"`
	Timestamp      string `json:"timestamp"`
}

type UserLoginEvent struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

type UserUpdatedEvent struct {
	UserID    string `json:"userId"`
	Timestamp string `json:"timestamp"`
}

// ── Publish methods ──────────────────────────────────────────────

func (p *UserEventPublisher) PublishUserCreated(userID, fullName, email, identityType, identityNumber string) error {
	return p.nats.Publish("auth.user.created", UserCreatedEvent{
		UserID:         userID,
		FullName:       fullName,
		Email:          email,
		IdentityType:   identityType,
		IdentityNumber: identityNumber,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
	})
}

func (p *UserEventPublisher) PublishUserLogin(userID, email string) error {
	return p.nats.Publish("auth.user.login", UserLoginEvent{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (p *UserEventPublisher) PublishUserUpdated(userID string) error {
	return p.nats.Publish("auth.user.updated", UserUpdatedEvent{
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
