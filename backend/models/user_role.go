package models

import "time"

type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	RoleID    uint      `gorm:"index;not null" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName - eksplisit agar GORM pakai "user_roles"
func (UserRole) TableName() string {
	return "user_roles"
}