package models

import (
    "time"
    // "gorm.io/gorm"
)

type User struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    UUID        string         `gorm:"unique;not null;index" json:"uuid"`
    Username    string         `gorm:"unique;not null;index" json:"username"`
    Email       string         `gorm:"unique;not null;index" json:"email"`
    Password    string         `gorm:"not null" json:"-"`
    FullName    string         `gorm:"not null" json:"full_name"`
    PhoneNumber string         `json:"phone_number"`
    Status      string         `gorm:"default:'active'" json:"status"`
    LastLogin   *time.Time     `json:"last_login"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    // DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    Roles       []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
    Documents   []Document     `gorm:"foreignKey:UserID" json:"documents,omitempty"`
}