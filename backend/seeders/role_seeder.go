package seeders

import (
    "log"
    "gorm.io/gorm"
    "atk-backend/models"
)

func SeedRoles(db *gorm.DB) {
    roles := []models.Role{
        {
            Name:        "super_admin",
            DisplayName: "Super Administrator",
            Description: "Full access to all features",
            Level:       5,
        },
        {
            Name:        "admin",
            DisplayName: "Administrator",
            Description: "Manage users, masters, and transactions",
            Level:       4,
        },
        {
            Name:        "manager",
            DisplayName: "Manager",
            Description: "Approve transactions and manage stock",
            Level:       3,
        },
        {
            Name:        "staff",
            DisplayName: "Staff",
            Description: "Create transactions and manage inventory",
            Level:       2,
        },
        {
            Name:        "viewer",
            DisplayName: "Viewer",
            Description: "Read-only access",
            Level:       1,
        },
    }
    
    for _, role := range roles {
        var existing models.Role
        result := db.Where("name = ?", role.Name).First(&existing)
        if result.Error != nil {
            db.Create(&role)
            log.Printf("  ✅ Created role: %s", role.Name)
        } else {
            log.Printf("  ⏭️  Role already exists: %s", role.Name)
        }
    }
}