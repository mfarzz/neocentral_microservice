package dto

import "time"

// ── Academic Year DTOs ─────────────────────────

type CreateAcademicYearRequest struct {
	Semester  string `json:"semester" validate:"required,oneof=ganjil genap"`
	Year      string `json:"year" validate:"required"`
	StartDate string `json:"startDate"` // format: 2024-08-01
	EndDate   string `json:"endDate"`
	IsActive  bool   `json:"isActive"`
}

type UpdateAcademicYearRequest struct {
	Semester  *string `json:"semester" validate:"omitempty,oneof=ganjil genap"`
	Year      *string `json:"year"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
	IsActive  *bool   `json:"isActive"`
}

type AcademicYearResponse struct {
	ID        string     `json:"id"`
	Semester  string     `json:"semester"`
	Year      *string    `json:"year"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ── Room DTOs ──────────────────────────────────

type CreateRoomRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location"`
	Capacity *int   `json:"capacity"`
}

type UpdateRoomRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
	Capacity *int    `json:"capacity"`
}

type RoomResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  *string   `json:"location"`
	Capacity  *int      `json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ── Science Group DTOs ────────────────────────

type CreateScienceGroupRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateScienceGroupRequest struct {
	Name *string `json:"name"`
}

type ScienceGroupResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ── Thesis Topic DTOs ────────────────────────

type CreateThesisTopicRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateThesisTopicRequest struct {
	Name *string `json:"name"`
}

type ThesisTopicResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ── Thesis Status DTOs ───────────────────────

type ThesisStatusResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
