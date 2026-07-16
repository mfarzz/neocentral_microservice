package repository

import (
	"context"

	"neocentral-go/master-data-service/internal/model"
)

// AcademicYearRepository defines the interface for academic year data access.
type AcademicYearRepository interface {
	FindAll(ctx context.Context) ([]model.AcademicYear, error)
	FindByID(ctx context.Context, id string) (*model.AcademicYear, error)
	FindActive(ctx context.Context) (*model.AcademicYear, error)
	Create(ctx context.Context, ay *model.AcademicYear) error
	Update(ctx context.Context, ay *model.AcademicYear) error
	Delete(ctx context.Context, id string) error
	DeactivateAll(ctx context.Context) error
}

// RoomRepository defines the interface for room data access.
type RoomRepository interface {
	FindAll(ctx context.Context) ([]model.Room, error)
	FindByID(ctx context.Context, id string) (*model.Room, error)
	Create(ctx context.Context, room *model.Room) error
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id string) error
}

// ScienceGroupRepository defines the interface for science group data access.
type ScienceGroupRepository interface {
	FindAll(ctx context.Context) ([]model.ScienceGroup, error)
	FindByID(ctx context.Context, id string) (*model.ScienceGroup, error)
	Create(ctx context.Context, sg *model.ScienceGroup) error
	Update(ctx context.Context, sg *model.ScienceGroup) error
	Delete(ctx context.Context, id string) error
}

// ThesisTopicRepository defines the interface for thesis topic data access.
type ThesisTopicRepository interface {
	FindAll(ctx context.Context) ([]model.ThesisTopic, error)
	FindByID(ctx context.Context, id string) (*model.ThesisTopic, error)
	Create(ctx context.Context, topic *model.ThesisTopic) error
	Update(ctx context.Context, topic *model.ThesisTopic) error
	Delete(ctx context.Context, id string) error
}

// ThesisStatusRepository defines the interface for thesis status data access.
type ThesisStatusRepository interface {
	FindAll(ctx context.Context) ([]model.ThesisStatus, error)
	FindByID(ctx context.Context, id string) (*model.ThesisStatus, error)
}
