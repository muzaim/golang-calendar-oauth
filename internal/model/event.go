package model

import (
	"time"
)

type Event struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"not null;index" json:"user_id"`
	GoogleEventID *string   `gorm:"type:varchar(255)" json:"google_event_id"`
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	Description   *string   `gorm:"type:text" json:"description"`
	Location      *string   `gorm:"type:varchar(255)" json:"location"`
	StartTime     time.Time `gorm:"type:datetime;not null" json:"start_time"`
	EndTime       time.Time `gorm:"type:datetime;not null" json:"end_time"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Event) TableName() string {
	return "events"
}
