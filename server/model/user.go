package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Name       string    `json:"name"`
	Password   string    `gorm:"not null" json:"password"`
	RoleID     uint      `gorm:"not null" json:"role_id"`
	Role       Role      `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE" json:"role"`
	LastAccess time.Time `json:"last_access"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
