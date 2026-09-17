package utils

import (
	"atk-backend/database"
	"gorm.io/gorm"
)

// GetDB returns database connection
func GetDB() *gorm.DB {
	return database.GetDB()
}