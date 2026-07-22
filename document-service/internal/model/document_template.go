package model

import "time"

type DocumentTemplate struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"type:varchar(191);uniqueIndex;not null" json:"name"`
	Type      string    `gorm:"type:varchar(191);default:'HTML'" json:"type"`
	Content   *string   `gorm:"type:longtext" json:"content"`
	FilePath  *string   `gorm:"type:varchar(255)" json:"filePath"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (DocumentTemplate) TableName() string {
	return "document_templates"
}
