package model

import "time"

// ThesisTopic represents the thesis_topics table.
type ThesisTopic struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ThesisTopic) TableName() string { return "thesis_topics" }

// ThesisStatus represents the thesis_status table.
type ThesisStatus struct {
	ID   string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

func (ThesisStatus) TableName() string { return "thesis_status" }
