package models

import (
    "time"
    "gorm.io/gorm"
)

type Master struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    UUID        string         `gorm:"unique;not null;index" json:"uuid"`
    Code        string         `gorm:"unique;not null;index" json:"code"`
    Name        string         `gorm:"not null;index" json:"name"`
    Category    string         `gorm:"not null;index" json:"category"`
    Description string         `json:"description"`
    Unit        string         `gorm:"not null" json:"unit"`
    MinStock    int            `gorm:"default:0" json:"min_stock"`
    MaxStock    int            `gorm:"default:0" json:"max_stock"`
    Price       float64        `gorm:"type:decimal(10,2)" json:"price"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    Stocks      []Stock        `gorm:"foreignKey:MasterID" json:"stocks,omitempty"`
}