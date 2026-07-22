package model

import "time"

type InternshipLogbook struct {
	ID                  string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InternshipID        string    `gorm:"type:varchar(36);not null" json:"internshipId"`
	ActivityDate        time.Time `gorm:"type:date;not null" json:"activityDate"`
	ActivityDescription string    `gorm:"type:text;not null" json:"activityDescription"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Internship *Internship `gorm:"foreignKey:InternshipID" json:"internship,omitempty"`
}

func (InternshipLogbook) TableName() string {
	return "internship_logbooks"
}
