package services

import (
	"errors"

	"atk-backend/models"
	"atk-backend/repositories"

	"github.com/google/uuid"
)

// ============================================================
// MASTER SERVICE
// ============================================================
type MasterService struct {
	masterRepo *repositories.MasterRepository
	stockRepo  *repositories.StockRepository
}

func NewMasterService() *MasterService {
	return &MasterService{
		masterRepo: repositories.NewMasterRepository(),
		stockRepo:  repositories.NewStockRepository(),
	}
}

// ============================================================
// REQUEST STRUCTS — camelCase untuk konsisten dengan frontend
// ============================================================
type CreateMasterRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	Description  string  `json:"description"`
	Unit         string  `json:"unit" binding:"required"`
	MinStock     int     `json:"minStock"`
	MaxStock     int     `json:"maxStock"`
	Price        float64 `json:"price"`
	InitialStock int     `json:"initialStock"`
}

type UpdateMasterRequest struct {
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
	Unit        string  `json:"unit" binding:"required"`
	MinStock    int     `json:"minStock"`
	MaxStock    int     `json:"maxStock"`
	Price       float64 `json:"price"`
}

// ============================================================
// CREATE
// ============================================================
func (s *MasterService) Create(req CreateMasterRequest) (*models.Master, error) {
	// Cek duplikat kode
	if _, err := s.masterRepo.FindByCode(req.Code); err == nil {
		return nil, errors.New("master code already exists")
	}

	master := &models.Master{
		UUID:        uuid.New().String(),
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		Unit:        req.Unit,
		MinStock:    req.MinStock,
		MaxStock:    req.MaxStock,
		Price:       req.Price,
	}

	if err := s.masterRepo.Create(master); err != nil {
		return nil, err
	}

	// Buat stock awal kalau initialStock > 0
	if req.InitialStock > 0 {
		stock := &models.Stock{
			MasterID: master.ID,
			Quantity: req.InitialStock,
			Location: "Gudang Utama",
			Status:   "available",
		}
		if err := s.stockRepo.Create(stock); err != nil {
			// Rollback? Tidak — user bisa tambah stock manual nanti
			// Log saja, tapi jangan gagalkan create master
			// (kalau mau transaksional, pakai DB transaction)
		}
	}

	return s.masterRepo.FindByID(master.ID)
}

// ============================================================
// GET BY ID
// ============================================================
func (s *MasterService) GetByID(id uint) (*models.Master, error) {
	return s.masterRepo.FindByID(id)
}

// ============================================================
// GET BY CODE
// ============================================================
func (s *MasterService) GetByCode(code string) (*models.Master, error) {
	return s.masterRepo.FindByCode(code)
}

// ============================================================
// GET ALL — dengan search + filter category
// ============================================================
func (s *MasterService) GetAll(page, limit int, search, category string) ([]models.Master, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.masterRepo.FindAll(offset, limit, search, category)
}

// ============================================================
// UPDATE
// ============================================================
func (s *MasterService) Update(id uint, req UpdateMasterRequest) (*models.Master, error) {
	master, err := s.masterRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("master not found")
	}

	master.Name = req.Name
	master.Category = req.Category
	master.Description = req.Description
	master.Unit = req.Unit
	master.MinStock = req.MinStock
	master.MaxStock = req.MaxStock
	master.Price = req.Price

	if err := s.masterRepo.Update(master); err != nil {
		return nil, err
	}
	return s.masterRepo.FindByID(id)
}

// ============================================================
// DELETE
// ============================================================
func (s *MasterService) Delete(id uint) error {
	return s.masterRepo.Delete(id)
}

// ============================================================
// BULK IMPORT
// ============================================================
func (s *MasterService) BulkImport(masters []models.Master) (int, error) {
	if len(masters) == 0 {
		return 0, errors.New("no data to import")
	}

	for i := range masters {
		masters[i].UUID = uuid.New().String()
	}

	if err := s.masterRepo.BulkInsertWithChunk(masters, 1000); err != nil {
		return 0, err
	}

	return len(masters), nil
}