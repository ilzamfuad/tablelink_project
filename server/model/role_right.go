package model

import "time"

type RoleRight struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Section   string    `gorm:"not null" json:"section"`
	Route     string    `gorm:"not null" json:"route"`
	RCreate   int       `gorm:"default:0" json:"r_create"`
	RRead     int       `gorm:"default:0" json:"r_read"`
	RUpdate   int       `gorm:"default:0" json:"r_update"`
	RDelete   int       `gorm:"default:0" json:"r_delete"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Role      Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role"`
}
