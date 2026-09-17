package main

import (
    "flag"
    "log"
    
    "atk-backend/config"
    "atk-backend/database"
    "atk-backend/migrations"
    "atk-backend/seeders"
    "atk-backend/routes"
    
    "gorm.io/gorm"
)

func main() {
    // Parse command line arguments
    migrateCmd := flag.String("migrate", "", "Run migration: up, down, fresh")
    seederCmd := flag.Bool("seeder", false, "Run seeders")
    seedersCmd := flag.Bool("seeders", false, "Run seeders (alias for -seeder)")
    flag.Parse()

    // Load config
    cfg := config.LoadConfig()
    
    // Initialize database
    database.InitDB(cfg)
    db := database.GetDB()

    // Handle migration commands
    if *migrateCmd != "" {
        handleMigration(db, *migrateCmd)
        
        // Jika ada flag seeder atau seeders, jalankan seeder setelah migration
        if *seederCmd || *seedersCmd {
            handleSeeder(db)
        }
        return
    }

    // Handle seeder only (without migration)
    if *seederCmd || *seedersCmd {
        handleSeeder(db)
        return
    }

    // Run server
    log.Println("🚀 Starting ATK Enterprise Server...")
    router := routes.SetupRouter()
    
    port := cfg.Port
    if port == "" {
        port = "8080"
    }
    
    log.Printf("✅ Server running on http://localhost:%s", port)
    log.Printf("📚 API Documentation available at http://localhost:%s/health", port)
    
    router.Run(":" + port)
}

func handleMigration(db *gorm.DB, command string) {
    switch command {
    case "up":
        log.Println("📦 Running migrations up...")
        migrations.RunMigrations(db)
        log.Println("✅ Migrations completed successfully")
    case "down":
        log.Println("📦 Running migrations down...")
        migrations.RollbackMigrations(db)
        log.Println("✅ Rollback completed successfully")
    case "fresh":
        log.Println("📦 Running fresh migrations...")
        migrations.DropAllTables(db)
        migrations.RunMigrations(db)
        log.Println("✅ Fresh migrations completed successfully")
    default:
        log.Fatalf("❌ Unknown migration command: %s\nUsage: go run main.go -migrate=up|down|fresh", command)
    }
}

func handleSeeder(db *gorm.DB) {
    log.Println("🌱 Running seeders...")
    seeders.RunSeeders(db)
    log.Println("✅ Seeders completed successfully")
}