package models

import (
	"time"
)

type Transaction struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TransactionCode string     `gorm:"unique;not null;index" json:"transactionCode"`
	Type            string     `gorm:"not null;index" json:"type"`
	Status          string     `gorm:"default:'draft';index" json:"status"`
	UserID          uint       `gorm:"not null;index" json:"userId"`
	ApprovedBy      *uint      `json:"approvedBy"`
	ApprovedAt      *time.Time `json:"approvedAt"`
	Notes           string     `json:"notes"`
	TotalItems      int        `json:"totalItems"`
	TotalValue      float64    `gorm:"type:decimal(10,2)" json:"totalValue"`
	TransactionDate time.Time  `gorm:"not null;index" json:"transactionDate"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`

	User      User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Approver  *User             `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Items     []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Documents []Document        `gorm:"foreignKey:TransactionID" json:"documents,omitempty"`
}