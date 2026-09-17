package seeders

import (
    "log"
    "gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) {
    log.Println("🌱 Running seeders...")
    
    // Seed roles
    log.Println("  📋 Seeding roles...")
    SeedRoles(db)
    
    // Seed users
    log.Println("  👤 Seeding users...")
    SeedUsers(db)
    
    // Seed masters
    log.Println("  📦 Seeding masters...")
    SeedMasters(db)
    
    log.Println("✅ Seeders completed successfully")
}