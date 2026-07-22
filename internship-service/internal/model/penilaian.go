package model

import "time"

type InternshipCPMK struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Code           string    `gorm:"type:varchar(255);not null" json:"code"`
	Name           string    `gorm:"type:text;not null" json:"name"`
	Weight         float64   `gorm:"not null" json:"weight"`
	AssessorType   string    `gorm:"type:varchar(20);not null" json:"assessorType"`
	AcademicYearID string    `gorm:"type:varchar(36);not null" json:"academicYearId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	Rubrics []InternshipAssessmentRubric `gorm:"foreignKey:CpmkID" json:"rubrics,omitempty"`
}

func (InternshipCPMK) TableName() string {
	return "internship_cpmks"
}

type InternshipAssessmentRubric struct {
	ID                     string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CpmkID                 string    `gorm:"type:varchar(36);not null;column:cpmk_id" json:"cpmkId"`
	LevelName              string    `gorm:"type:varchar(255);not null" json:"levelName"`
	RubricLevelDescription string    `gorm:"type:text;not null" json:"rubricLevelDescription"`
	MinScore               float64   `gorm:"not null" json:"minScore"`
	MaxScore               float64   `gorm:"not null" json:"maxScore"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

func (InternshipAssessmentRubric) TableName() string {
	return "internship_assessment_rubrics"
}

type InternshipLecturerScore struct {
	InternshipID   string    `gorm:"primaryKey;type:varchar(36)" json:"internshipId"`
	ChosenRubricID string    `gorm:"primaryKey;type:varchar(36)" json:"chosenRubricId"`
	Score          float64   `gorm:"not null" json:"score"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (InternshipLecturerScore) TableName() string {
	return "internship_lecturer_scores"
}

type InternshipFieldScore struct {
	InternshipID   string    `gorm:"primaryKey;type:varchar(36)" json:"internshipId"`
	ChosenRubricID string    `gorm:"primaryKey;type:varchar(36)" json:"chosenRubricId"`
	Score          float64   `gorm:"not null" json:"score"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (InternshipFieldScore) TableName() string {
	return "internship_field_scores"
}

type FieldAssessmentToken struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InternshipID string     `gorm:"type:varchar(36);not null;column:internship_id" json:"internshipId"`
	Token        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Pin          *string    `gorm:"type:varchar(6)" json:"pin"`
	ExpiresAt    time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	IsUsed       bool       `gorm:"default:false;column:is_used" json:"isUsed"`
	UsedAt       *time.Time `gorm:"column:used_at" json:"usedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func (FieldAssessmentToken) TableName() string {
	return "field_assessment_tokens"
}
