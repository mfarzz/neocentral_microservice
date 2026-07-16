package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"neocentral-go/auth-service/internal/domain"
)

// MockUserRepository is a mock implementation of repository.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByIdentityNumber(ctx context.Context, identityNumber string) (*domain.User, error) {
	args := m.Called(ctx, identityNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	args := m.Called(ctx, id, fields)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.UserHasRole, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.UserHasRole), args.Error(1)
}

func (m *MockUserRepository) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
	args := m.Called(ctx, userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUsersByRole(ctx context.Context, roleName string, page, pageSize int) ([]domain.User, int64, error) {
	args := m.Called(ctx, roleName, page, pageSize)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) BatchGetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]domain.User), args.Error(1)
}
