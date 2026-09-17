package repositories

import (
	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository() *StockRepository {
	return &StockRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (r *StockRepository) Create(stock *models.Stock) error {
	return r.db.Create(stock).Error
}

// ============================================================
// FIND BY ID
// ============================================================
func (r *StockRepository) FindByID(id uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.Preload("Master").First(&stock, id).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// ============================================================
// FIND BY MASTER ID
// ============================================================
func (r *StockRepository) FindByMasterID(masterID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.Preload("Master").Where("master_id = ?", masterID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// ============================================================
// EXISTS BY MASTER ID
// ============================================================
func (r *StockRepository) ExistsByMasterID(masterID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Stock{}).Where("master_id = ?", masterID).Count(&count).Error
	return count > 0, err
}

// ============================================================
// UPDATE
// ============================================================
func (r *StockRepository) Update(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

// ============================================================
// DELETE
// ============================================================
func (r *StockRepository) Delete(id uint) error {
	return r.db.Delete(&models.Stock{}, id).Error
}

// ============================================================
// UPDATE STOCK QUANTITY (atomic increment/decrement)
// ============================================================
func (r *StockRepository) UpdateStockQuantity(masterID uint, quantity int) error {
	return r.db.Model(&models.Stock{}).
		Where("master_id = ?", masterID).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

// ============================================================
// FIND ALL — basic pagination (backward compat)
// ============================================================
func (r *StockRepository) FindAll(offset, limit int) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	query := r.db.Model(&models.Stock{}).Preload("Master")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&stocks).Error

	return stocks, total, err
}

// ============================================================
// FIND ALL WITH MASTER — pagination + search + status filter
// Digunakan oleh StockHandler.GetAll (versi baru)
// ============================================================
func (r *StockRepository) FindAllWithMaster(
	offset, limit int,
	search, status string,
) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	// Base query dengan Preload relasi Master
	query := r.db.Model(&models.Stock{}).
		Preload("Master").
		Joins("LEFT JOIN masters ON masters.id = stocks.master_id")

	// Filter search — cari di kode master, nama master, atau lokasi stock
	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"masters.code LIKE ? OR masters.name LIKE ? OR stocks.location LIKE ?",
			like, like, like,
		)
	}

	// Filter status
	if status != "" {
		query = query.Where("stocks.status = ?", status)
	}

	// Count total (sebelum offset/limit)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data
	err := query.
		Order("stocks.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&stocks).Error

	return stocks, total, err
}

// ============================================================
// FIND LOW STOCK
// ============================================================
func (r *StockRepository) FindLowStock(minStock int) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Preload("Master").
		Where("quantity <= ?", minStock).
		Order("quantity ASC").
		Find(&stocks).Error
	return stocks, err
}

// ============================================================
// BULK UPDATE
// ============================================================
func (r *StockRepository) BulkUpdate(stocks []models.Stock) error {
	for _, stock := range stocks {
		err := r.db.Model(&models.Stock{}).
			Where("master_id = ?", stock.MasterID).
			Updates(map[string]interface{}{
				"quantity":     stock.Quantity,
				"last_updated": stock.LastUpdated,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}