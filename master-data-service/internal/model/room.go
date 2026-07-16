package model

import "time"

// Room represents the rooms table.
type Room struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Location  *string   `gorm:"type:varchar(255)" json:"location"`
	Capacity  *int      `json:"capacity"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Room) TableName() string { return "rooms" }
