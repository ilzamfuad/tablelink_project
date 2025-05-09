package model

import (
	"time"
)

type Token struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccessToken  string    `gorm:"type:varchar(500);unique_index;not_null" json:"access_token"`
	RefreshToken string    `gorm:"type:varchar(500);unique_index;not_null" json:"refresh_token"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	ExpiredAt    time.Time `gorm:"type:timestamptz;not_null" json:"expired_at"`
}
