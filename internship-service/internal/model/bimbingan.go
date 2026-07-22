package model

import "time"

type InternshipGuidanceQuestion struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WeekNumber     int       `gorm:"not null" json:"weekNumber"`
	QuestionText   string    `gorm:"type:text;not null" json:"questionText"`
	OrderIndex     int       `gorm:"default:0" json:"orderIndex"`
	AcademicYearID string    `gorm:"type:varchar(36);not null" json:"academicYearId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type GuidanceCriteriaInputType string

const (
	CriteriaInputEvaluation GuidanceCriteriaInputType = "EVALUATION"
	CriteriaInputText       GuidanceCriteriaInputType = "TEXT"
)

type InternshipGuidanceLecturerCriteria struct {
	ID             string                    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CriteriaName   string                    `gorm:"type:varchar(255);not null" json:"criteriaName"`
	WeekNumber     int                       `gorm:"not null" json:"weekNumber"`
	InputType      GuidanceCriteriaInputType `gorm:"type:enum('EVALUATION', 'TEXT');not null" json:"inputType"`
	OrderIndex     int                       `gorm:"default:0" json:"orderIndex"`
	AcademicYearID string                    `gorm:"type:varchar(36);not null" json:"academicYearId"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`

	Options []InternshipGuidanceLecturerCriteriaOption `gorm:"foreignKey:CriteriaID" json:"options,omitempty"`
}

type InternshipGuidanceLecturerCriteriaOption struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CriteriaID string    `gorm:"type:varchar(36);not null" json:"criteriaId"`
	OptionText string    `gorm:"type:varchar(255);not null" json:"optionText"`
	OrderIndex int       `gorm:"default:0" json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GuidanceSessionStatus string

const (
	SessionSubmitted GuidanceSessionStatus = "SUBMITTED"
	SessionLate      GuidanceSessionStatus = "LATE"
	SessionApproved  GuidanceSessionStatus = "APPROVED"
)

type InternshipGuidanceSession struct {
	ID             string                `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InternshipID   string                `gorm:"type:varchar(36);not null" json:"internshipId"`
	WeekNumber     int                   `gorm:"not null" json:"weekNumber"`
	Status         GuidanceSessionStatus `gorm:"type:enum('SUBMITTED', 'LATE', 'APPROVED');default:'SUBMITTED'" json:"status"`
	SubmissionDate *time.Time            `gorm:"type:date" json:"submissionDate"`
	ApprovedAt     *time.Time            `json:"approvedAt"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`

	StudentAnswers  []InternshipGuidanceStudentAnswer  `gorm:"foreignKey:GuidanceSessionID" json:"studentAnswers,omitempty"`
	LecturerAnswers []InternshipGuidanceLecturerAnswer `gorm:"foreignKey:GuidanceSessionID" json:"lecturerAnswers,omitempty"`
}

type InternshipGuidanceStudentAnswer struct {
	GuidanceSessionID string    `gorm:"primaryKey;type:varchar(36)" json:"guidanceSessionId"`
	QuestionID        string    `gorm:"primaryKey;type:varchar(36)" json:"questionId"`
	WeekNumber        int       `gorm:"primaryKey" json:"weekNumber"`
	AnswerText        string    `gorm:"type:text;not null" json:"answerText"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type InternshipGuidanceLecturerAnswer struct {
	GuidanceSessionID string    `gorm:"primaryKey;type:varchar(36)" json:"guidanceSessionId"`
	CriteriaID        string    `gorm:"primaryKey;type:varchar(36)" json:"criteriaId"`
	WeekNumber        int       `gorm:"primaryKey" json:"weekNumber"`
	EvaluationValue   *string   `gorm:"type:varchar(255)" json:"evaluationValue"`
	AnswerText        *string   `gorm:"type:text" json:"answerText"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
