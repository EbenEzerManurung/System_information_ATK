package repositories

import (
	"time"

	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// ============================================================
// FIND BY ID — dengan preload lengkap
// ============================================================
func (r *TransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.
		Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master").
		Preload("Documents").
		First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ============================================================
// FIND BY CODE
// ============================================================
func (r *TransactionRepository) FindByCode(code string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.
		Where("transaction_code = ?", code).
		Preload("User").
		Preload("Items").
		Preload("Items.Master").
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ============================================================
// FIND ALL (backward compat)
// ============================================================
func (r *TransactionRepository) FindAll(
	offset, limit int,
	status, transType string,
	userID uint,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).
		Preload("User").
		Preload("Approver")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if transType != "" {
		query = query.Where("type = ?", transType)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

// ============================================================
// FIND ALL WITH SEARCH — untuk handler GetAll & ExportToExcel
// ============================================================
func (r *TransactionRepository) FindAllWithSearch(
	offset, limit int,
	search, status, transType string,
	userID uint,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).
		Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"transaction_code LIKE ? OR notes LIKE ?",
			like, like,
		)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if transType != "" {
		query = query.Where("type = ?", transType)
	}

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

// ============================================================
// UPDATE
// ============================================================
func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// ============================================================
// UPDATE STATUS — pakai time.Now() supaya konsisten
// ============================================================
func (r *TransactionRepository) UpdateStatus(
	id uint,
	status string,
	approvedBy *uint,
) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if approvedBy != nil {
		now := time.Now()
		updates["approved_by"] = *approvedBy
		updates["approved_at"] = now
	}

	return r.db.
		Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// ============================================================
// DELETE
// ============================================================
func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// ============================================================
// FIND BY USER ID
// ============================================================
func (r *TransactionRepository) FindByUserID(
	userID uint,
	offset, limit int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

// ============================================================
// COUNT BY STATUS — untuk dashboard summary
// ============================================================
func (r *TransactionRepository) CountByStatus(status string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Transaction{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

// ============================================================
// BULK UPDATE STATUS
// ============================================================
func (r *TransactionRepository) BulkUpdateStatus(ids []uint, status string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.
		Model(&models.Transaction{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// ============================================================
// ✅ DELETE ITEMS — hapus semua item transaksi (untuk update)
// ============================================================
func (r *TransactionRepository) DeleteItems(transactionID uint) error {
	return r.db.
		Where("transaction_id = ?", transactionID).
		Delete(&models.TransactionItem{}).Error
}

// ============================================================
// ✅ REPLACE ITEMS — hapus item lama + create items baru (transactional)
// ============================================================
func (r *TransactionRepository) ReplaceItems(
	transactionID uint,
	items []models.TransactionItem,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Hapus semua item lama
		if err := tx.
			Where("transaction_id = ?", transactionID).
			Delete(&models.TransactionItem{}).Error; err != nil {
			return err
		}

		// 2. Kalau ada items baru, set TransactionID dan create
		if len(items) > 0 {
			for i := range items {
				items[i].TransactionID = transactionID
				items[i].ID = 0 // reset ID supaya auto-increment
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ============================================================
// ✅ UPDATE WITH ITEMS — update transaksi + replace items (transactional)
// Method lengkap untuk update all-in-one
// ============================================================
func (r *TransactionRepository) UpdateWithItems(
	transaction *models.Transaction,
	items []models.TransactionItem,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update transaksi utama
		if err := tx.Save(transaction).Error; err != nil {
			return err
		}

		// 2. Hapus items lama
		if err := tx.
			Where("transaction_id = ?", transaction.ID).
			Delete(&models.TransactionItem{}).Error; err != nil {
			return err
		}

		// 3. Create items baru
		if len(items) > 0 {
			for i := range items {
				items[i].TransactionID = transaction.ID
				items[i].ID = 0
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}