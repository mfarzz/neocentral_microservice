package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"neocentral-go/auth-service/internal/domain"
	"neocentral-go/auth-service/internal/service"
	"neocentral-go/auth-service/tests/mocks"
	"neocentral-go/pkg/apperror"
)

func TestProfileService_GetProfile(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewProfileService(mockRepo)
	ctx := context.Background()

	t.Run("Success - User Found", func(t *testing.T) {
		userID := "user-1"
		mockUser := &domain.User{
			ID:             userID,
			FullName:       "John Doe",
			IdentityNumber: "12345",
			IdentityType:   domain.IdentityNIM,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.On("FindByID", ctx, userID).Return(mockUser, nil).Once()

		resp, err := svc.GetProfile(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "John Doe", resp.FullName)
		assert.Equal(t, "12345", resp.IdentityNumber)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userID := "user-2"

		mockRepo.On("FindByID", ctx, userID).Return((*domain.User)(nil), nil).Once()

		resp, err := svc.GetProfile(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		
		appErr, ok := err.(*apperror.AppError)
		assert.True(t, ok)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "User not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})
}
