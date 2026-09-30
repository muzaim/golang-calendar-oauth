package model

import (
	"time"
)

type User struct {
	ID                    int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string     `gorm:"type:varchar(255);not null" json:"name"`
	Email                 string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash          *string    `gorm:"type:varchar(255)" json:"-"`
	Provider              string     `gorm:"type:varchar(50);not null;default:'email'" json:"provider"`
	ProviderID            *string    `gorm:"type:varchar(255)" json:"provider_id"`
	Avatar                *string    `gorm:"type:text" json:"avatar"`
	GoogleAccessToken     *string    `gorm:"type:text" json:"-"`
	GoogleRefreshToken    *string    `gorm:"type:text" json:"-"`
	GoogleTokenExpiresAt *time.Time `gorm:"type:datetime" json:"-"`
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Events        []Event        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (User) TableName() string {
	return "users"
}
