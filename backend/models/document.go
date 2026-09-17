package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UUID          string     `gorm:"unique;not null;index" json:"uuid"`
	DocumentCode  string     `gorm:"unique;not null;index" json:"documentCode"`
	Type          string     `gorm:"not null;index" json:"type"`
	TransactionID *uint      `gorm:"index" json:"transactionId"`
	UserID        uint       `gorm:"not null;index" json:"userId"`
	QRCode        string     `gorm:"type:text" json:"qrCode,omitempty"`
	QRCodeImage   string     `gorm:"type:longtext" json:"qrCodeImage,omitempty"`
	Signature     string     `gorm:"type:longtext" json:"signature,omitempty"`
	SignedAt      *time.Time `json:"signedAt,omitempty"`
	Status        string     `gorm:"default:'pending';index" json:"status"`
	FilePath      string     `gorm:"size:500" json:"filePath,omitempty"`
	Metadata      string     `gorm:"type:text" json:"metadata,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`

	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ✅ Auto-generate UUID sebelum create
func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.UUID == "" {
		d.UUID = uuid.New().String()
	}
	return nil
}