package model

import "time"

// ScienceGroup represents the science_groups table (Kelompok Keahlian).
type ScienceGroup struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ScienceGroup) TableName() string { return "science_groups" }
