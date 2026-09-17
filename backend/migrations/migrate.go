package migrations

import (
    "log"
    "gorm.io/gorm"
    "atk-backend/models"
)

func RunMigrations(db *gorm.DB) {
    log.Println("📦 Running migrations...")
    
    // Create tables with proper order for foreign keys
    err := db.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.UserRole{},
        &models.Master{},
        &models.Stock{},
        &models.Transaction{},
        &models.TransactionItem{},
        &models.Document{},
    )
    
    if err != nil {
        log.Fatal("❌ Migration failed:", err)
    }
    
    log.Println("✅ Migrations completed successfully")
}

func RollbackMigrations(db *gorm.DB) {
    log.Println("📦 Rolling back migrations...")
    
    // Drop tables in reverse order
    err := db.Migrator().DropTable(
        &models.Document{},
        &models.TransactionItem{},
        &models.Transaction{},
        &models.Stock{},
        &models.Master{},
        &models.UserRole{},
        &models.Role{},
        &models.User{},
    )
    
    if err != nil {
        log.Fatal("❌ Rollback failed:", err)
    }
    
    log.Println("✅ Rollback completed successfully")
}

func DropAllTables(db *gorm.DB) {
    log.Println("📦 Dropping all tables...")
    
    err := db.Migrator().DropTable(
        &models.Document{},
        &models.TransactionItem{},
        &models.Transaction{},
        &models.Stock{},
        &models.Master{},
        &models.UserRole{},
        &models.Role{},
        &models.User{},
    )
    
    if err != nil {
        log.Fatal("❌ Drop tables failed:", err)
    }
    
    log.Println("✅ All tables dropped successfully")
}