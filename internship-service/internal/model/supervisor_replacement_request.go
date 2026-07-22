package model

import "time"

type SupervisorReplacementStatus string

const (
	ReplacementPending  SupervisorReplacementStatus = "PENDING"
	ReplacementApproved SupervisorReplacementStatus = "APPROVED"
	ReplacementRejected SupervisorReplacementStatus = "REJECTED"
)

type SupervisorReplacementRequest struct {
	ID                string                      `gorm:"primaryKey;type:varchar(36)" json:"id"`
	LetterID          string                      `gorm:"type:varchar(36);not null;column:letter_id" json:"letterId"`
	InternshipID      string                      `gorm:"type:varchar(36);not null;column:internship_id" json:"internshipId"`
	OldSupervisorID   string                      `gorm:"type:varchar(36);not null;column:old_supervisor_id" json:"oldSupervisorId"`
	NewSupervisorID   string                      `gorm:"type:varchar(36);not null;column:new_supervisor_id" json:"newSupervisorId"`
	Reason            string                      `gorm:"type:text;not null" json:"reason"`
	Status            SupervisorReplacementStatus `gorm:"type:enum('PENDING', 'APPROVED', 'REJECTED');default:'PENDING'" json:"status"`
	RequestedByID     string                      `gorm:"type:varchar(36);not null;column:requested_by_id" json:"requestedById"`
	ApprovedByID      *string                     `gorm:"type:varchar(36);column:approved_by_id" json:"approvedById"`
	RequestedAt       time.Time                   `gorm:"column:requested_at;autoCreateTime" json:"requestedAt"`
	ResolvedAt        *time.Time                  `gorm:"column:resolved_at" json:"resolvedAt"`
	RejectionNotes    *string                     `gorm:"type:text;column:rejection_notes" json:"rejectionNotes"`
	CreatedAt         time.Time                   `json:"createdAt"`
	UpdatedAt         time.Time                   `json:"updatedAt"`
}

func (SupervisorReplacementRequest) TableName() string {
	return "supervisor_replacement_requests"
}
