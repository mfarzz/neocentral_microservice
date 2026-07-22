package model

import "time"

type SeminarStatus string

const (
	SeminarRequested SeminarStatus = "REQUESTED"
	SeminarApproved  SeminarStatus = "APPROVED"
	SeminarRejected  SeminarStatus = "REJECTED"
	SeminarCompleted SeminarStatus = "COMPLETED"
	SeminarFailed    SeminarStatus = "FAILED"
)

type AudienceStatus string

const (
	AudiencePending   AudienceStatus = "PENDING"
	AudienceValidated AudienceStatus = "VALIDATED"
	AudienceRejected  AudienceStatus = "REJECTED"
)

type InternshipSeminar struct {
	ID                     string        `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InternshipID           string        `gorm:"type:varchar(36);not null" json:"internshipId"`
	RoomID                 string        `gorm:"type:varchar(36);not null" json:"roomId"`
	SeminarDate            time.Time     `gorm:"type:date;not null" json:"seminarDate"`
	StartTime              string        `gorm:"type:time;not null" json:"startTime"`
	EndTime                string        `gorm:"type:time;not null" json:"endTime"`
	LinkMeeting            string        `gorm:"type:varchar(255)" json:"linkMeeting"`
	ModeratorStudentID     string        `gorm:"type:varchar(36);not null" json:"moderatorStudentId"`
	Status                 SeminarStatus `gorm:"type:enum('REQUESTED', 'APPROVED', 'REJECTED', 'COMPLETED', 'FAILED');default:'REQUESTED'" json:"status"`
	ApprovedBy             string        `gorm:"type:varchar(36)" json:"approvedBy"`
	SupervisorNotes        string        `gorm:"type:text" json:"supervisorNotes"`
	BeritaAcaraDocumentID  string        `gorm:"type:varchar(36)" json:"beritaAcaraDocumentId"`
	CreatedAt              time.Time     `json:"createdAt"`
	UpdatedAt              time.Time     `json:"updatedAt"`
}

func (InternshipSeminar) TableName() string {
	return "internship_seminars"
}

type InternshipSeminarAudience struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	SeminarID   string         `gorm:"type:varchar(36);not null" json:"seminarId"`
	StudentID   string         `gorm:"type:varchar(36);not null" json:"studentId"`
	Status      AudienceStatus `gorm:"type:enum('PENDING', 'VALIDATED', 'REJECTED');default:'PENDING'" json:"status"`
	ValidatedBy string         `gorm:"type:varchar(36)" json:"validatedBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func (InternshipSeminarAudience) TableName() string {
	return "internship_seminar_audiences"
}
