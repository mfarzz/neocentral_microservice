package model

import "time"

// AcademicYear represents the academic_years table.
type AcademicYear struct {
	ID        string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Semester  string     `gorm:"type:enum('ganjil','genap');default:'ganjil'" json:"semester"`
	Year      *string    `gorm:"type:varchar(20)" json:"year"`
	StartDate *time.Time `gorm:"column:start_date" json:"startDate"`
	EndDate   *time.Time `gorm:"column:end_date" json:"endDate"`
	IsActive  bool       `gorm:"column:is_active;default:false" json:"isActive"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AcademicYear) TableName() string { return "academic_years" }
