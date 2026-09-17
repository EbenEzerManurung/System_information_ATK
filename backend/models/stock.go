package models

import (
    "time"
)

type Stock struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    MasterID    uint      `gorm:"not null;index" json:"master_id"`
    Quantity    int       `gorm:"not null" json:"quantity"`
    Location    string    `json:"location"`
    Status      string    `gorm:"default:'available';index" json:"status"`
    LastUpdated time.Time `json:"last_updated"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Master      Master    `gorm:"foreignKey:MasterID" json:"master,omitempty"`
}