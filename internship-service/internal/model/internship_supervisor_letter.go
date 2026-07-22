package model

import "time"

type SupervisorLetterStatus string

const (
	SupLetterActive     SupervisorLetterStatus = "ACTIVE"
	SupLetterSuperseded SupervisorLetterStatus = "SUPERSEDED"
)

type InternshipSupervisorLetter struct {
	ID               string                 `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DocumentNumber   string                 `gorm:"type:varchar(255);unique;not null" json:"documentNumber"`
	DateIssued       time.Time              `gorm:"type:date;not null" json:"dateIssued"`
	StartDate        time.Time              `gorm:"type:date;not null" json:"startDate"`
	EndDate          time.Time              `gorm:"type:date;not null" json:"endDate"`
	SupervisorID     string                 `gorm:"type:varchar(36);not null" json:"supervisorId"`
	DocumentID       *string                `gorm:"type:varchar(36)" json:"documentId"`
	SignedByID       *string                `gorm:"type:varchar(36)" json:"signedById"`
	SignedAsRoleID   *string                `gorm:"type:varchar(36)" json:"signedAsRoleId"`
	Status           SupervisorLetterStatus `gorm:"type:enum('ACTIVE', 'SUPERSEDED');default:'ACTIVE'" json:"status"`
	SupersededAt     *time.Time             `json:"supersededAt"`
	SupersededReason *string                `gorm:"type:text" json:"supersededReason"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
}

func (InternshipSupervisorLetter) TableName() string {
	return "internship_supervisor_letters"
}
