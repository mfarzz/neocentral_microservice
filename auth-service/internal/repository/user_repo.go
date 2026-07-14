package repository

import (
	"context"

	"neocentral-go/auth-service/internal/domain"
)

// UserRepository defines the data-access contract for User entities.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	GetUserRoles(ctx context.Context, userID string) ([]domain.UserHasRole, error)
	HasRole(ctx context.Context, userID, roleName string) (bool, error)
	GetUsersByRole(ctx context.Context, roleName string, page, pageSize int) ([]domain.User, int64, error)
	BatchGetByIDs(ctx context.Context, ids []string) ([]domain.User, error)
}
