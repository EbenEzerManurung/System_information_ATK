package models

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null;index" json:"name"`
	DisplayName string    `gorm:"not null" json:"displayName"`
	Description string    `json:"description"`
	Level       int       `gorm:"default:1" json:"level"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Users []User `gorm:"many2many:user_roles;" json:"users,omitempty"`
}