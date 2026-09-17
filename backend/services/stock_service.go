package services

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"atk-backend/models"
	"atk-backend/repositories"

	"github.com/xuri/excelize/v2"
)

// ============================================================
// STOCK SERVICE
// ============================================================
type StockService struct {
	stockRepo  *repositories.StockRepository
	masterRepo *repositories.MasterRepository
}

func NewStockService() *StockService {
	return &StockService{
		stockRepo:  repositories.NewStockRepository(),
		masterRepo: repositories.NewMasterRepository(),
	}
}

// ============================================================
// REQUEST STRUCTS
// ============================================================
type CreateStockRequest struct {
	MasterID uint   `json:"masterId" binding:"required"`
	Quantity int    `json:"quantity"`
	Location string `json:"location" binding:"required"`
	Status   string `json:"status"`
}

type UpdateStockRequest struct {
	Location string `json:"location"`
	Status   string `json:"status"`
	Quantity *int   `json:"quantity,omitempty"`
}

type ImportStockResult struct {
	Imported int      `json:"imported"`
	Updated  int      `json:"updated"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

// ============================================================
// CREATE
// ============================================================
func (s *StockService) Create(req CreateStockRequest) (*models.Stock, error) {
	// Cek master exist
	master, err := s.masterRepo.FindByID(req.MasterID)
	if err != nil || master == nil {
		return nil, errors.New("master tidak ditemukan")
	}

	// Cek apakah stock untuk master ini sudah ada
	exists, err := s.stockRepo.ExistsByMasterID(req.MasterID)
	if err == nil && exists {
		return nil, errors.New("stock untuk master ini sudah ada. Gunakan Edit untuk update")
	}

	status := req.Status
	if status == "" {
		status = "available"
	}

	stock := &models.Stock{
		MasterID: req.MasterID,
		Quantity: req.Quantity,
		Location: req.Location,
		Status:   status,
	}

	if err := s.stockRepo.Create(stock); err != nil {
		return nil, err
	}
	return s.stockRepo.FindByID(stock.ID)
}

// ============================================================
// GET BY ID
// ============================================================
func (s *StockService) GetByID(id uint) (*models.Stock, error) {
	return s.stockRepo.FindByID(id)
}

// ============================================================
// GET BY MASTER ID
// ============================================================
func (s *StockService) GetByMasterID(masterID uint) (*models.Stock, error) {
	return s.stockRepo.FindByMasterID(masterID)
}

// ============================================================
// GET ALL — basic pagination (backward compat)
// ============================================================
func (s *StockService) GetAll(page, limit int) ([]models.Stock, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.stockRepo.FindAll(offset, limit)
}

// ============================================================
// GET ALL WITH MASTER — pagination + search + status filter
// ============================================================
func (s *StockService) GetAllWithMaster(page, size int, search, status string) ([]models.Stock, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size
	return s.stockRepo.FindAllWithMaster(offset, size, search, status)
}

// ============================================================
// UPDATE
// ============================================================
func (s *StockService) Update(id uint, req UpdateStockRequest) (*models.Stock, error) {
	stock, err := s.stockRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("stock tidak ditemukan")
	}

	if req.Location != "" {
		stock.Location = req.Location
	}
	if req.Status != "" {
		stock.Status = req.Status
	}
	if req.Quantity != nil {
		stock.Quantity = *req.Quantity
	}

	if err := s.stockRepo.Update(stock); err != nil {
		return nil, err
	}
	return s.stockRepo.FindByID(id)
}

// ============================================================
// DELETE
// ============================================================
func (s *StockService) Delete(id uint) error {
	return s.stockRepo.Delete(id)
}

// ============================================================
// GET LOW STOCK
// ============================================================
func (s *StockService) GetLowStock(minStock int) ([]models.Stock, error) {
	return s.stockRepo.FindLowStock(minStock)
}

// ============================================================
// GET STOCK SUMMARY
// ✅ FIX: konversi int64 → int agar bisa dioperasikan
// ============================================================
func (s *StockService) GetStockSummary() (interface{}, error) {
	all, total, err := s.stockRepo.FindAll(0, 10000)
	if err != nil {
		return nil, err
	}

	totalQty := 0
	lowCount := 0
	outCount := 0

	for _, st := range all {
		totalQty += st.Quantity
		if st.Quantity == 0 {
			outCount++
		} else if st.Quantity <= 10 {
			lowCount++
		}
	}

	// ✅ Konversi total dari int64 ke int
	totalInt := int(total)
	available := totalInt - lowCount - outCount
	if available < 0 {
		available = 0
	}

	return map[string]interface{}{
		"totalItems": totalInt,
		"totalQty":   totalQty,
		"lowStock":   lowCount,
		"outOfStock": outCount,
		"available":  available,
	}, nil
}

// ============================================================
// ADJUST STOCK
// ============================================================
func (s *StockService) AdjustStock(masterID uint, quantity int, notes string) (*models.Stock, error) {
	stock, err := s.stockRepo.FindByMasterID(masterID)
	if err != nil {
		return nil, errors.New("stock tidak ditemukan untuk master ini")
	}

	newQty := stock.Quantity + quantity
	if newQty < 0 {
		return nil, errors.New("quantity tidak boleh negatif")
	}

	stock.Quantity = newQty

	// Auto-update status berdasarkan quantity
	if newQty == 0 {
		stock.Status = "out"
	} else if newQty <= 10 {
		stock.Status = "low"
	} else {
		stock.Status = "available"
	}

	if err := s.stockRepo.Update(stock); err != nil {
		return nil, err
	}

	return s.stockRepo.FindByMasterID(masterID)
}

// ============================================================
// IMPORT EXCEL
// Format kolom (baris 1 = header):
// A: Master Code | B: Quantity | C: Location | D: Status
// ============================================================
func (s *StockService) ImportFromExcel(fh *multipart.FileHeader) (*ImportStockResult, error) {
	if fh == nil {
		return nil, errors.New("file tidak valid")
	}
	if fh.Size > 5*1024*1024 {
		return nil, errors.New("ukuran file maksimal 5MB")
	}

	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("file bukan format excel valid: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("file excel kosong")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("gagal baca baris: %w", err)
	}
	if len(rows) < 2 {
		return nil, errors.New("file excel tidak memiliki data")
	}

	result := &ImportStockResult{Errors: []string{}}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		if len(row) < 2 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: kolom tidak lengkap (minimal kode master & quantity)", i+1))
			continue
		}

		masterCode := strings.TrimSpace(row[0])
		quantity := 0
		fmt.Sscanf(strings.TrimSpace(row[1]), "%d", &quantity)

		location := ""
		if len(row) >= 3 {
			location = strings.TrimSpace(row[2])
		}
		if location == "" {
			location = "Gudang Utama"
		}

		status := "available"
		if len(row) >= 4 && strings.TrimSpace(row[3]) != "" {
			status = strings.TrimSpace(row[3])
		}

		// Cari master by code
		master, err := s.masterRepo.FindByCode(masterCode)
		if err != nil || master == nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: master code '%s' tidak ditemukan", i+1, masterCode))
			continue
		}

		// Cek apakah stock sudah ada
		existing, err := s.stockRepo.FindByMasterID(master.ID)
		if err == nil && existing != nil && existing.ID > 0 {
			// Update
			existing.Quantity = quantity
			existing.Location = location
			existing.Status = status
			if err := s.stockRepo.Update(existing); err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("Baris %d: gagal update - %s", i+1, err.Error()))
				continue
			}
			result.Updated++
		} else {
			// Create
			stock := &models.Stock{
				MasterID: master.ID,
				Quantity: quantity,
				Location: location,
				Status:   status,
			}
			if err := s.stockRepo.Create(stock); err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("Baris %d: gagal simpan - %s", i+1, err.Error()))
				continue
			}
			result.Imported++
		}
	}

	return result, nil
}

// ============================================================
// EXPORT EXCEL
// ============================================================
func (s *StockService) ExportToExcel(search, status string) (*bytes.Buffer, error) {
	stocks, _, err := s.stockRepo.FindAllWithMaster(0, 10000, search, status)
	if err != nil {
		return nil, fmt.Errorf("gagal query stock: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Stock"
	idx, _ := f.NewSheet(sheet)
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2563EB"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "1E40AF", Style: 1},
			{Type: "right", Color: "1E40AF", Style: 1},
			{Type: "top", Color: "1E40AF", Style: 1},
			{Type: "bottom", Color: "1E40AF", Style: 1},
		},
	})

	headers := []string{"No", "Kode", "Nama Item", "Kategori", "Qty", "Satuan", "Lokasi", "Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, st := range stocks {
		row := i + 2
		masterCode, masterName, masterCategory, masterUnit := "", "", "", ""
		if st.Master.ID > 0 {
			masterCode = st.Master.Code
			masterName = st.Master.Name
			masterCategory = st.Master.Category
			masterUnit = st.Master.Unit
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), masterCode)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), masterName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), masterCategory)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), st.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), masterUnit)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), st.Location)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), st.Status)
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 10)
	f.SetColWidth(sheet, "F", "F", 10)
	f.SetColWidth(sheet, "G", "G", 20)
	f.SetColWidth(sheet, "H", "H", 15)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: 1,
		TopLeftCell: "A2", ActivePane: "bottomLeft",
	})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("gagal generate file: %w", err)
	}
	return buf, nil
}