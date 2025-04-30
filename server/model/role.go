package model

import "time"

type Role struct {
	ID        uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string      `json:"name"`
	RoleRight []RoleRight `gorm:"foreignKey:RoleID;references:ID" json:"role_rights"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
