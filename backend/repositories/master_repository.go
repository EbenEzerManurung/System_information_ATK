package repositories

import (
	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type MasterRepository struct {
	db *gorm.DB
}

func NewMasterRepository() *MasterRepository {
	return &MasterRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (r *MasterRepository) Create(master *models.Master) error {
	return r.db.Create(master).Error
}

// ============================================================
// FIND BY ID
// ============================================================
func (r *MasterRepository) FindByID(id uint) (*models.Master, error) {
	var master models.Master
	err := r.db.Preload("Stocks").First(&master, id).Error
	if err != nil {
		return nil, err
	}
	return &master, nil
}

// ============================================================
// FIND BY CODE
// ============================================================
func (r *MasterRepository) FindByCode(code string) (*models.Master, error) {
	var master models.Master
	err := r.db.Where("code = ?", code).First(&master).Error
	if err != nil {
		return nil, err
	}
	return &master, nil
}

// ============================================================
// FIND ALL — pagination + search + filter category
// ============================================================
func (r *MasterRepository) FindAll(
	offset, limit int,
	search, category string,
) ([]models.Master, int64, error) {
	var masters []models.Master
	var total int64

	query := r.db.Model(&models.Master{}).Preload("Stocks")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"code LIKE ? OR name LIKE ? OR description LIKE ?",
			like, like, like,
		)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// ✅ Cek error Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("code ASC").
		Offset(offset).
		Limit(limit).
		Find(&masters).Error

	return masters, total, err
}

// ============================================================
// FIND ALL WITH SEARCH — untuk export (tanpa pagination)
// ============================================================
func (r *MasterRepository) FindAllWithSearch(
	offset, limit int,
	search, category string,
) ([]models.Master, error) {
	var masters []models.Master

	query := r.db.Model(&models.Master{}).Preload("Stocks")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"code LIKE ? OR name LIKE ? OR description LIKE ?",
			like, like, like,
		)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.
		Order("code ASC").
		Offset(offset).
		Limit(limit).
		Find(&masters).Error

	return masters, err
}

// ============================================================
// UPDATE
// ============================================================
func (r *MasterRepository) Update(master *models.Master) error {
	return r.db.Save(master).Error
}

// ============================================================
// DELETE
// ============================================================
func (r *MasterRepository) Delete(id uint) error {
	return r.db.Delete(&models.Master{}, id).Error
}

// ============================================================
// FIND ALL FOR EXPORT
// ============================================================
func (r *MasterRepository) FindAllForExport() ([]models.Master, error) {
	var masters []models.Master
	err := r.db.Order("code ASC").Find(&masters).Error
	return masters, err
}

// ============================================================
// BULK INSERT
// ============================================================
func (r *MasterRepository) BulkInsert(masters []models.Master) error {
	return r.db.CreateInBatches(masters, 1000).Error
}

// ============================================================
// BULK INSERT WITH CHUNK
// ============================================================
func (r *MasterRepository) BulkInsertWithChunk(
	masters []models.Master,
	chunkSize int,
) error {
	if chunkSize <= 0 {
		chunkSize = 1000
	}
	return r.db.CreateInBatches(masters, chunkSize).Error
}

// ============================================================
// COUNT ALL
// ============================================================
func (r *MasterRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Master{}).Count(&count).Error
	return count, err
}