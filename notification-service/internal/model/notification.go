package model

import "time"

type Notification struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index;not null" json:"userId"`
	Title     *string   `gorm:"type:varchar(255)" json:"title"`
	Message   *string   `gorm:"type:text" json:"message"`
	IsRead    bool      `gorm:"default:false;not null" json:"isRead"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
}

func (Notification) TableName() string {
	return "notifications"
}
