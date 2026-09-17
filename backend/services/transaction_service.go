package services

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"atk-backend/models"
	"atk-backend/repositories"

	"github.com/xuri/excelize/v2"
)

// ============================================================
// TRANSACTION SERVICE
// ============================================================
type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	masterRepo      *repositories.MasterRepository
	stockRepo       *repositories.StockRepository
	documentRepo    *repositories.DocumentRepository
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		transactionRepo: repositories.NewTransactionRepository(),
		masterRepo:      repositories.NewMasterRepository(),
		stockRepo:       repositories.NewStockRepository(),
		documentRepo:    repositories.NewDocumentRepository(),
	}
}

// ============================================================
// REQUEST STRUCTS — camelCase
// ============================================================
type CreateTransactionRequest struct {
	Type   string                   `json:"type" binding:"required,oneof=IN OUT RETURN ADJUSTMENT"`
	Notes  string                   `json:"notes"`
	Status string                   `json:"status"`
	Items  []TransactionItemRequest `json:"items" binding:"required,min=1"`
}

type TransactionItemRequest struct {
	MasterID uint   `json:"masterId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Notes    string `json:"notes"`
}

type UpdateTransactionRequest struct {
	Notes  string `json:"notes"`
	Status string `json:"status"`
}

// Struct untuk approve dengan signature
type ApproveTransactionRequest struct {
	Signature     string `json:"signature"`
	SignatureData string `json:"signatureData"`
	Notes         string `json:"notes"`
}

// Struct untuk reject
type RejectTransactionRequest struct {
	Signature string `json:"signature"`
	Notes     string `json:"notes"`
	Reason    string `json:"reason"`
}

type ImportTransactionResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

// ============================================================
// CREATE
// ============================================================
func (s *TransactionService) Create(req CreateTransactionRequest, userID uint) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("transaction must have at least one item")
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "pending" {
		status = "draft"
	}

	code := fmt.Sprintf("TRX-%s-%04d",
		time.Now().Format("20060102"),
		time.Now().UnixNano()%10000)

	transaction := &models.Transaction{
		TransactionCode: code,
		Type:            req.Type,
		Status:          status,
		UserID:          userID,
		Notes:           req.Notes,
		TransactionDate: time.Now(),
	}

	var totalItems int
	var totalValue float64
	var items []models.TransactionItem

	for _, itemReq := range req.Items {
		master, err := s.masterRepo.FindByID(itemReq.MasterID)
		if err != nil {
			return nil, fmt.Errorf("master with ID %d not found", itemReq.MasterID)
		}

		if req.Type == "OUT" || req.Type == "ADJUSTMENT" {
			stock, err := s.stockRepo.FindByMasterID(itemReq.MasterID)
			if err != nil {
				return nil, fmt.Errorf("stock untuk %s tidak ditemukan", master.Name)
			}
			if stock.Quantity < itemReq.Quantity {
				return nil, fmt.Errorf("stock tidak cukup untuk %s (tersedia: %d, diminta: %d)",
					master.Name, stock.Quantity, itemReq.Quantity)
			}
		}

		subTotal := float64(itemReq.Quantity) * master.Price
		totalItems += itemReq.Quantity
		totalValue += subTotal

		item := models.TransactionItem{
			MasterID: itemReq.MasterID,
			Quantity: itemReq.Quantity,
			Price:    master.Price,
			SubTotal: subTotal,
			Notes:    itemReq.Notes,
		}
		items = append(items, item)
	}

	transaction.TotalItems = totalItems
	transaction.TotalValue = totalValue
	transaction.Items = items

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	// Update stock
	for _, item := range transaction.Items {
		quantity := item.Quantity
		if req.Type == "OUT" {
			quantity = -quantity
		} else if req.Type == "ADJUSTMENT" {
			stock, _ := s.stockRepo.FindByMasterID(item.MasterID)
			if stock != nil {
				quantity = item.Quantity - stock.Quantity
			}
		}
		s.stockRepo.UpdateStockQuantity(item.MasterID, quantity)
	}

	return s.transactionRepo.FindByID(transaction.ID)
}

// ============================================================
// GET BY ID
// ============================================================
func (s *TransactionService) GetByID(id uint) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id)
}

// ============================================================
// GET BY CODE
// ============================================================
func (s *TransactionService) GetByCode(code string) (*models.Transaction, error) {
	return s.transactionRepo.FindByCode(code)
}

// ============================================================
// GET ALL (backward compat)
// ============================================================
func (s *TransactionService) GetAll(
	page, limit int,
	status, transType string,
	userID uint,
) ([]models.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.transactionRepo.FindAll(offset, limit, status, transType, userID)
}

// ============================================================
// GET ALL WITH SEARCH
// ============================================================
func (s *TransactionService) GetAllWithSearch(
	page, size int,
	search, status, transType string,
	userID uint,
) ([]models.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size
	return s.transactionRepo.FindAllWithSearch(
		offset, size, search, status, transType, userID,
	)
}

// ============================================================
// UPDATE
// ============================================================
func (s *TransactionService) Update(id uint, req UpdateTransactionRequest) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction.Status != "draft" && transaction.Status != "pending" {
		return nil, errors.New("hanya transaksi draft atau pending yang dapat diupdate")
	}

	if req.Notes != "" {
		transaction.Notes = req.Notes
	}
	if req.Status != "" && (req.Status == "draft" || req.Status == "pending") {
		transaction.Status = req.Status
	}

	if err := s.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}
	return s.transactionRepo.FindByID(id)
}

// ============================================================
// DELETE
// ============================================================
func (s *TransactionService) Delete(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != "draft" && transaction.Status != "cancelled" {
		return errors.New("hanya transaksi draft atau cancelled yang dapat dihapus")
	}

	return s.transactionRepo.Delete(id)
}

// ============================================================
// APPROVE — ✅ FIX: pakai field Signature (bukan SignatureData)
// ============================================================
func (s *TransactionService) Approve(
	id uint,
	approverID uint,
	req ApproveTransactionRequest,
) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	// Terima status draft ATAU pending
	if transaction.Status != "draft" && transaction.Status != "pending" {
		return nil, fmt.Errorf(
			"transaksi dengan status '%s' tidak dapat disetujui (hanya draft/pending)",
			transaction.Status,
		)
	}

	// Ambil signature dari salah satu field
	signature := req.Signature
	if signature == "" {
		signature = req.SignatureData
	}
	if signature == "" {
		return nil, errors.New("tanda tangan wajib diisi")
	}

	log.Printf("[Transaction.Approve] ID=%d, ApproverID=%d, SigLen=%d",
		id, approverID, len(signature))

	// Update status transaksi → approved
	if err := s.transactionRepo.UpdateStatus(id, "approved", &approverID); err != nil {
		return nil, err
	}

	// Update notes jika ada
	if req.Notes != "" {
		updated, _ := s.transactionRepo.FindByID(id)
		updated.Notes = req.Notes
		s.transactionRepo.Update(updated)
	}

	// ✅ FIX: Pakai CreateSignature (method khusus) — field Signature, bukan SignatureData
	docService := NewDocumentService()
	doc, err := docService.CreateSignature(id, approverID, signature, req.Notes)
	if err != nil {
		log.Printf("[Transaction.Approve] ERROR: gagal simpan signature: %v", err)
		return nil, fmt.Errorf("gagal simpan signature: %w", err)
	}

	log.Printf("[Transaction.Approve] Signature saved: DocID=%d, Code=%s",
		doc.ID, doc.DocumentCode)

	return s.transactionRepo.FindByID(id)
}

// ============================================================
// REJECT
// ============================================================
func (s *TransactionService) Reject(
	id uint,
	approverID uint,
	req RejectTransactionRequest,
) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaction.Status != "draft" && transaction.Status != "pending" {
		return nil, fmt.Errorf(
			"transaksi dengan status '%s' tidak dapat ditolak",
			transaction.Status,
		)
	}

	if err := s.transactionRepo.UpdateStatus(id, "rejected", &approverID); err != nil {
		return nil, err
	}

	reason := req.Reason
	if reason == "" {
		reason = req.Notes
	}
	if reason != "" {
		updated, _ := s.transactionRepo.FindByID(id)
		updated.Notes = reason
		s.transactionRepo.Update(updated)
	}

	return s.transactionRepo.FindByID(id)
}

// ============================================================
// SUBMIT
// ============================================================
func (s *TransactionService) Submit(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != "draft" {
		return errors.New("only draft transactions can be submitted")
	}

	return s.transactionRepo.UpdateStatus(id, "pending", nil)
}

// ============================================================
// CANCEL
// ============================================================
func (s *TransactionService) Cancel(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status == "completed" || transaction.Status == "cancelled" {
		return errors.New("transaction cannot be cancelled")
	}

	if transaction.Type == "OUT" && transaction.Status != "cancelled" {
		for _, item := range transaction.Items {
			s.stockRepo.UpdateStockQuantity(item.MasterID, item.Quantity)
		}
	}

	return s.transactionRepo.UpdateStatus(id, "cancelled", nil)
}

// ============================================================
// COMPLETE
// ============================================================
func (s *TransactionService) Complete(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != "approved" {
		return errors.New("transaction must be approved first")
	}

	return s.transactionRepo.UpdateStatus(id, "completed", nil)
}

// ============================================================
// EXPORT EXCEL
// ============================================================
func (s *TransactionService) ExportToExcel(
	search, status, transType string,
) (*bytes.Buffer, error) {
	transactions, _, err := s.transactionRepo.FindAllWithSearch(
		0, 10000, search, status, transType, 0,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal query transaksi: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Transactions"
	idx, _ := f.NewSheet(sheet)
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")

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

	headers := []string{
		"No", "Kode", "Tipe", "Status", "Tanggal",
		"Total Item", "Total Value", "Catatan",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, t := range transactions {
		row := i + 2
		dateStr := ""
		if !t.TransactionDate.IsZero() {
			dateStr = t.TransactionDate.Format("2006-01-02 15:04:05")
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.TransactionCode)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.Type)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), dateStr)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), t.TotalItems)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), t.TotalValue)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), t.Notes)
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 22)
	f.SetColWidth(sheet, "C", "C", 12)
	f.SetColWidth(sheet, "D", "D", 12)
	f.SetColWidth(sheet, "E", "E", 20)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 15)
	f.SetColWidth(sheet, "H", "H", 30)

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

// ============================================================
// IMPORT EXCEL
// ============================================================
func (s *TransactionService) ImportFromExcel(
	fh *multipart.FileHeader,
	userID uint,
) (*ImportTransactionResult, error) {
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
		return nil, fmt.Errorf("file bukan excel valid: %w", err)
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

	result := &ImportTransactionResult{Errors: []string{}}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue
		}
		if strings.TrimSpace(row[0]) == "" || strings.TrimSpace(row[1]) == "" {
			continue
		}

		transType := strings.ToUpper(strings.TrimSpace(row[0]))
		masterCode := strings.TrimSpace(row[1])
		quantity, _ := strconv.Atoi(strings.TrimSpace(row[2]))
		notes := ""
		if len(row) >= 4 {
			notes = strings.TrimSpace(row[3])
		}

		if transType != "IN" && transType != "OUT" {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: tipe '%s' tidak valid (IN/OUT)", i+1, transType))
			continue
		}
		if quantity <= 0 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: quantity harus > 0", i+1))
			continue
		}

		master, err := s.masterRepo.FindByCode(masterCode)
		if err != nil || master == nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: master code '%s' tidak ditemukan", i+1, masterCode))
			continue
		}

		req := CreateTransactionRequest{
			Type:   transType,
			Notes:  notes,
			Status: "draft",
			Items: []TransactionItemRequest{
				{
					MasterID: master.ID,
					Quantity: quantity,
					Notes:    notes,
				},
			},
		}

		if _, err := s.Create(req, userID); err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: gagal simpan - %s", i+1, err.Error()))
			continue
		}

		result.Imported++
	}

	return result, nil
}