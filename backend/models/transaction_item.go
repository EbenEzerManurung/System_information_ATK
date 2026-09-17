package models

import (
	"time"
)

type TransactionItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null;index" json:"transactionId"`
	MasterID      uint      `gorm:"not null;index" json:"masterId"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"type:decimal(10,2)" json:"price"`
	SubTotal      float64   `gorm:"type:decimal(10,2)" json:"subTotal"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`

	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Master      Master      `gorm:"foreignKey:MasterID" json:"master,omitempty"`
}