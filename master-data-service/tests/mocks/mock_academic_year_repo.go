package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"neocentral-go/master-data-service/internal/model"
)

type MockAcademicYearRepo struct {
	mock.Mock
}

func (m *MockAcademicYearRepo) FindAll(ctx context.Context) ([]model.AcademicYear, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.AcademicYear), args.Error(1)
}

func (m *MockAcademicYearRepo) FindByID(ctx context.Context, id string) (*model.AcademicYear, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AcademicYear), args.Error(1)
}

func (m *MockAcademicYearRepo) FindActive(ctx context.Context) (*model.AcademicYear, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AcademicYear), args.Error(1)
}

func (m *MockAcademicYearRepo) Create(ctx context.Context, ay *model.AcademicYear) error {
	args := m.Called(ctx, ay)
	return args.Error(0)
}

func (m *MockAcademicYearRepo) Update(ctx context.Context, ay *model.AcademicYear) error {
	args := m.Called(ctx, ay)
	return args.Error(0)
}

func (m *MockAcademicYearRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAcademicYearRepo) DeactivateAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
