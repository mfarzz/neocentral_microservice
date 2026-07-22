package model

import "time"

type InternshipHoliday struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	HolidayDate time.Time `gorm:"type:date;not null;unique" json:"holidayDate"`
	Name        *string   `gorm:"type:varchar(255)" json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (InternshipHoliday) TableName() string {
	return "internship_holidays"
}
