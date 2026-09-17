package seeders

import (
    "log"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "atk-backend/models"
)

func SeedUsers(db *gorm.DB) {
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    
    // Get roles
    var superAdminRole, adminRole, staffRole models.Role
    db.Where("name = ?", "super_admin").First(&superAdminRole)
    db.Where("name = ?", "admin").First(&adminRole)
    db.Where("name = ?", "staff").First(&staffRole)
    
    // Create Super Admin
    superAdmin := models.User{
        UUID:        uuid.New().String(),
        Username:    "superadmin",
        Email:       "superadmin@atk.com",
        Password:    string(hashedPassword),
        FullName:    "Super Administrator",
        PhoneNumber: "081234567890",
        Status:      "active",
    }
    
    var existing models.User
    result := db.Where("username = ?", "superadmin").First(&existing)
    if result.Error != nil {
        db.Create(&superAdmin)
        db.Create(&models.UserRole{
            UserID: superAdmin.ID,
            RoleID: superAdminRole.ID,
        })
        log.Printf("  ✅ Created user: superadmin")
    } else {
        log.Printf("  ⏭️  User already exists: superadmin")
    }
    
    // Create Admin
    admin := models.User{
        UUID:        uuid.New().String(),
        Username:    "admin",
        Email:       "admin@atk.com",
        Password:    string(hashedPassword),
        FullName:    "Administrator",
        PhoneNumber: "081234567891",
        Status:      "active",
    }
    
    result = db.Where("username = ?", "admin").First(&existing)
    if result.Error != nil {
        db.Create(&admin)
        db.Create(&models.UserRole{
            UserID: admin.ID,
            RoleID: adminRole.ID,
        })
        log.Printf("  ✅ Created user: admin")
    } else {
        log.Printf("  ⏭️  User already exists: admin")
    }
    
    // Create Staff
    staff := models.User{
        UUID:        uuid.New().String(),
        Username:    "staff",
        Email:       "staff@atk.com",
        Password:    string(hashedPassword),
        FullName:    "Staff User",
        PhoneNumber: "081234567892",
        Status:      "active",
    }
    
    result = db.Where("username = ?", "staff").First(&existing)
    if result.Error != nil {
        db.Create(&staff)
        db.Create(&models.UserRole{
            UserID: staff.ID,
            RoleID: staffRole.ID,
        })
        log.Printf("  ✅ Created user: staff")
    } else {
        log.Printf("  ⏭️  User already exists: staff")
    }
}