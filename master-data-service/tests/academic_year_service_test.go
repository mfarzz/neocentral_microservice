package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/model"
	"neocentral-go/master-data-service/internal/service"
	"neocentral-go/master-data-service/tests/mocks"
	"neocentral-go/pkg/apperror"
)

func TestAcademicYearService_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockAcademicYearRepo)
	svc := service.NewAcademicYearService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockData := []model.AcademicYear{
			{ID: "1", Semester: "ganjil"},
			{ID: "2", Semester: "genap"},
		}

		mockRepo.On("FindAll", ctx).Return(mockData, nil).Once()

		resp, err := svc.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "1", resp[0].ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestAcademicYearService_GetActive(t *testing.T) {
	mockRepo := new(mocks.MockAcademicYearRepo)
	svc := service.NewAcademicYearService(mockRepo)
	ctx := context.Background()

	t.Run("Success - Found Active", func(t *testing.T) {
		mockData := &model.AcademicYear{ID: "active-id", IsActive: true}

		mockRepo.On("FindActive", ctx).Return(mockData, nil).Once()

		resp, err := svc.GetActive(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "active-id", resp.ID)
		assert.True(t, resp.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - No Active Found", func(t *testing.T) {
		mockRepo.On("FindActive", ctx).Return((*model.AcademicYear)(nil), nil).Once()

		resp, err := svc.GetActive(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)

		appErr, ok := err.(*apperror.AppError)
		assert.True(t, ok)
		assert.Equal(t, 404, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestAcademicYearService_Create(t *testing.T) {
	mockRepo := new(mocks.MockAcademicYearRepo)
	svc := service.NewAcademicYearService(mockRepo)
	ctx := context.Background()

	t.Run("Success - Create Inactive", func(t *testing.T) {
		req := dto.CreateAcademicYearRequest{
			Semester:  "ganjil",
			Year:      "2025",
			StartDate: "2025-08-01",
			EndDate:   "2026-01-31",
			IsActive:  false,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AcademicYear")).Return(nil).Once()

		resp, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "ganjil", resp.Semester)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Create Active (Deactivates Others)", func(t *testing.T) {
		req := dto.CreateAcademicYearRequest{
			Semester: "genap",
			IsActive: true,
		}

		mockRepo.On("DeactivateAll", ctx).Return(nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AcademicYear")).Return(nil).Once()

		resp, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid Date Format", func(t *testing.T) {
		req := dto.CreateAcademicYearRequest{
			StartDate: "08-01-2025", // Invalid format, expected YYYY-MM-DD
		}

		resp, err := svc.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		
		appErr, ok := err.(*apperror.AppError)
		assert.True(t, ok)
		assert.Equal(t, 400, appErr.Code)
	})
}
