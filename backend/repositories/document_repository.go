package repositories

import (
	"time"

	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (r *DocumentRepository) Create(document *models.Document) error {
	return r.db.Create(document).Error
}

// ============================================================
// FIND BY ID
// ============================================================
func (r *DocumentRepository) FindByID(id uint) (*models.Document, error) {
	var document models.Document
	err := r.db.
		Preload("User").
		Preload("Transaction").
		First(&document, id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// ============================================================
// FIND BY CODE
// ============================================================
func (r *DocumentRepository) FindByCode(code string) (*models.Document, error) {
	var document models.Document
	err := r.db.
		Where("document_code = ?", code).
		Preload("User").
		Preload("Transaction").
		First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// ============================================================
// FIND BY QR CODE
// ============================================================
func (r *DocumentRepository) FindByQRCode(qrCode string) (*models.Document, error) {
	var document models.Document
	err := r.db.
		Where("qr_code = ?", qrCode).
		Preload("User").
		Preload("Transaction").
		First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// ============================================================
// FIND ALL (backward compat)
// ============================================================
func (r *DocumentRepository) FindAll(
	offset, limit int,
	status, docType string,
) ([]models.Document, int64, error) {
	var documents []models.Document
	var total int64

	query := r.db.Model(&models.Document{}).
		Preload("User").
		Preload("Transaction")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if docType != "" {
		query = query.Where("type = ?", docType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&documents).Error

	return documents, total, err
}

// ============================================================
// FIND ALL WITH FILTER — ✅ method yang dipanggil service
// ============================================================
func (r *DocumentRepository) FindAllWithFilter(
	offset, limit int,
	status, docType string,
) ([]models.Document, int64, error) {
	return r.FindAll(offset, limit, status, docType)
}

// ============================================================
// UPDATE
// ============================================================
func (r *DocumentRepository) Update(document *models.Document) error {
	return r.db.Save(document).Error
}

// ============================================================
// UPDATE SIGNATURE — ✅ pakai time.Now() bukan gorm.Expr
// ============================================================
func (r *DocumentRepository) UpdateSignature(id uint, signature string) error {
	now := time.Now()
	return r.db.
		Model(&models.Document{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"signature": signature,
			"signed_at": now,
			"status":    "signed",
		}).Error
}

// ============================================================
// DELETE
// ============================================================
func (r *DocumentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Document{}, id).Error
}

// ============================================================
// FIND BY TRANSACTION ID
// ============================================================
func (r *DocumentRepository) FindByTransactionID(
	transactionID uint,
) ([]models.Document, error) {
	var documents []models.Document
	err := r.db.
		Where("transaction_id = ?", transactionID).
		Preload("User").
		Order("id DESC").
		Find(&documents).Error
	return documents, err
}

// ============================================================
// COUNT BY TYPE — untuk dashboard summary
// ============================================================
func (r *DocumentRepository) CountByType(docType string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Document{})
	if docType != "" {
		query = query.Where("type = ?", docType)
	}
	err := query.Count(&count).Error
	return count, err
}