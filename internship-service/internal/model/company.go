package model

import "time"

type CompanyStatus string

const (
	CompanySave      CompanyStatus = "save"
	CompanyBlacklist CompanyStatus = "blacklist"
	CompanyDiajukan  CompanyStatus = "diajukan"
)

type Company struct {
	ID             string        `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CompanyName    string        `gorm:"type:varchar(255);not null" json:"companyName"`
	CompanyAddress string        `gorm:"type:text;not null" json:"companyAddress"`
	Alasan         *string       `gorm:"type:text" json:"alasan"`
	Status         CompanyStatus `gorm:"type:enum('save', 'blacklist', 'diajukan');default:'save'" json:"status"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}

func (Company) TableName() string {
	return "companies"
}
