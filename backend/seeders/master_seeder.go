package seeders

import (
    "log"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "atk-backend/models"
)

func SeedMasters(db *gorm.DB) {
    masters := []models.Master{
        {
            UUID:        uuid.New().String(),
            Code:        "AT001",
            Name:        "Ballpoint Pen 0.5mm",
            Category:    "alat tulis",
            Description: "Ballpoint pen dengan ujung 0.5mm, tinta hitam",
            Unit:        "pcs",
            MinStock:    100,
            MaxStock:    1000,
            Price:       2500,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT002",
            Name:        "Pencil 2B",
            Category:    "alat tulis",
            Description: "Pensil 2B dengan penghapus",
            Unit:        "pcs",
            MinStock:    50,
            MaxStock:    500,
            Price:       1500,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT003",
            Name:        "A4 Paper 70gsm",
            Category:    "kertas",
            Description: "Kertas A4 70gsm, 500 sheets per rim",
            Unit:        "rim",
            MinStock:    10,
            MaxStock:    100,
            Price:       45000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT004",
            Name:        "Stapler Mini",
            Category:    "perlengkapan",
            Description: "Stapler mini dengan 100 staples",
            Unit:        "pcs",
            MinStock:    20,
            MaxStock:    200,
            Price:       12000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT005",
            Name:        "Correction Tape",
            Category:    "alat tulis",
            Description: "Correction tape roll 5m x 5mm",
            Unit:        "pcs",
            MinStock:    30,
            MaxStock:    300,
            Price:       8000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT006",
            Name:        "Highlighter Neon",
            Category:    "alat tulis",
            Description: "Highlighter neon warna kuning",
            Unit:        "pcs",
            MinStock:    25,
            MaxStock:    250,
            Price:       5000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT007",
            Name:        "Post-it Notes 76x76",
            Category:    "perlengkapan",
            Description: "Post-it notes ukuran 76x76mm, kuning",
            Unit:        "pack",
            MinStock:    15,
            MaxStock:    150,
            Price:       15000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT008",
            Name:        "USB Flash Drive 32GB",
            Category:    "elektronik",
            Description: "USB 3.0 32GB flash drive",
            Unit:        "pcs",
            MinStock:    5,
            MaxStock:    50,
            Price:       150000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT009",
            Name:        "Wireless Mouse",
            Category:    "elektronik",
            Description: "Wireless mouse 2.4GHz",
            Unit:        "pcs",
            MinStock:    3,
            MaxStock:    30,
            Price:       120000,
        },
        {
            UUID:        uuid.New().String(),
            Code:        "AT010",
            Name:        "Notebook A5 60 Sheets",
            Category:    "kertas",
            Description: "Notebook A5, 60 sheets, lined paper",
            Unit:        "pcs",
            MinStock:    20,
            MaxStock:    200,
            Price:       25000,
        },
    }
    
    for _, master := range masters {
        var existing models.Master
        result := db.Where("code = ?", master.Code).First(&existing)
        if result.Error != nil {
            db.Create(&master)
            
            // Create initial stock
            stock := models.Stock{
                MasterID:    master.ID,
                Quantity:    500,
                Location:    "Gudang Utama",
                Status:      "available",
                LastUpdated: master.CreatedAt,
            }
            db.Create(&stock)
            log.Printf("  ✅ Created master: %s", master.Code)
        }
    }
}