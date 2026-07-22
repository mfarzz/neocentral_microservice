package model

import "time"

type Document struct {
	ID             string        `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID         *string       `gorm:"type:varchar(36);index" json:"userId"`
	DocumentTypeID *string       `gorm:"type:varchar(36);index" json:"documentTypeId"`
	S3Bucket       string        `gorm:"type:varchar(191);not null" json:"s3Bucket"`
	S3ObjectName   string        `gorm:"type:varchar(255);not null" json:"s3ObjectName"`
	FileName       *string       `gorm:"type:varchar(191)" json:"fileName"`
	FileHash       *string       `gorm:"type:varchar(255)" json:"fileHash"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`

	DocumentType   *DocumentType `gorm:"foreignKey:DocumentTypeID" json:"documentType,omitempty"`
}

func (Document) TableName() string {
	return "documents"
}
