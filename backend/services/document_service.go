package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"atk-backend/database"
	"atk-backend/models"
	"atk-backend/repositories"
	"atk-backend/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================================
// DOCUMENT SERVICE
// ============================================================
type DocumentService struct {
	documentRepo    *repositories.DocumentRepository
	transactionRepo *repositories.TransactionRepository
	userRepo        *repositories.UserRepository
	db              *gorm.DB
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo:    repositories.NewDocumentRepository(),
		transactionRepo: repositories.NewTransactionRepository(),
		userRepo:        repositories.NewUserRepository(),
		db:              database.GetDB(),
	}
}

// ============================================================
// REQUEST STRUCTS
// ============================================================
type CreateDocumentRequest struct {
	Type          string `json:"type" binding:"required"`
	TransactionID *uint  `json:"transactionId"`
	UserID        uint   `json:"userId"`
	Signature     string `json:"signature"`
	Metadata      string `json:"metadata"`
	Status        string `json:"status"`
}

type ScanQRResponse struct {
	Document    models.Document     `json:"document"`
	Transaction *models.Transaction `json:"transaction,omitempty"`
	User        models.User         `json:"user"`
}

// ============================================================
// CREATE
// ============================================================
func (s *DocumentService) Create(req CreateDocumentRequest) (*models.Document, error) {
	code := fmt.Sprintf("DOC-%s-%04d",
		time.Now().Format("20060102"),
		time.Now().UnixNano()%10000)

	status := req.Status
	if status == "" {
		status = "pending"
	}

	document := &models.Document{
		UUID:          uuid.New().String(),
		DocumentCode:  code,
		Type:          req.Type,
		TransactionID: req.TransactionID,
		UserID:        req.UserID,
		Signature:     req.Signature,
		Status:        status,
		Metadata:      req.Metadata,
	}

	if req.Signature != "" {
		now := time.Now()
		document.SignedAt = &now
		document.Status = "signed"
	}

	qrData := map[string]interface{}{
		"documentCode": code,
		"type":         req.Type,
		"userId":       req.UserID,
		"createdAt":    time.Now().Format(time.RFC3339),
	}

	if req.TransactionID != nil {
		transaction, _ := s.transactionRepo.FindByID(*req.TransactionID)
		if transaction != nil {
			qrData["transactionCode"] = transaction.TransactionCode
			qrData["transactionId"] = *req.TransactionID
		}
	}

	qrJSON, _ := json.Marshal(qrData)
	document.QRCode = string(qrJSON)

	if qrImage, err := utils.GenerateQRCode(string(qrJSON)); err == nil {
		document.QRCodeImage = qrImage
	} else {
		log.Printf("[Document.Create] Warning: gagal generate QR: %v", err)
	}

	if err := s.documentRepo.Create(document); err != nil {
		return nil, fmt.Errorf("gagal simpan document: %w", err)
	}

	log.Printf("[Document.Create] ID=%d, Code=%s", document.ID, document.DocumentCode)
	return s.documentRepo.FindByID(document.ID)
}

// ============================================================
// CREATE SIGNATURE
// ============================================================
func (s *DocumentService) CreateSignature(
	transactionID uint,
	userID uint,
	signature string,
	notes string,
) (*models.Document, error) {
	if signature == "" {
		return nil, errors.New("signature wajib diisi")
	}

	code := fmt.Sprintf("SIG-%s-%04d",
		time.Now().Format("20060102"),
		time.Now().UnixNano()%10000)

	now := time.Now()
	metadata := fmt.Sprintf(`{"action":"approve","notes":%q}`, notes)

	document := &models.Document{
		UUID:          uuid.New().String(),
		DocumentCode:  code,
		Type:          "SIGNATURE",
		TransactionID: &transactionID,
		UserID:        userID,
		Signature:     signature,
		SignedAt:      &now,
		Status:        "signed",
		Metadata:      metadata,
	}

	qrData := map[string]interface{}{
		"documentCode":  code,
		"type":          "SIGNATURE",
		"transactionId": transactionID,
		"userId":        userID,
		"signedAt":      now.Format(time.RFC3339),
	}
	qrJSON, _ := json.Marshal(qrData)
	document.QRCode = string(qrJSON)

	if qrImage, err := utils.GenerateQRCode(string(qrJSON)); err == nil {
		document.QRCodeImage = qrImage
	}

	if err := s.documentRepo.Create(document); err != nil {
		return nil, fmt.Errorf("gagal simpan signature: %w", err)
	}

	log.Printf("[Document.CreateSignature] ID=%d, Code=%s, SigLen=%d",
		document.ID, document.DocumentCode, len(signature))

	return s.documentRepo.FindByID(document.ID)
}

// ============================================================
// GET BY ID / CODE / ALL
// ============================================================
func (s *DocumentService) GetByID(id uint) (*models.Document, error) {
	return s.documentRepo.FindByID(id)
}

func (s *DocumentService) GetByCode(code string) (*models.Document, error) {
	return s.documentRepo.FindByCode(code)
}

func (s *DocumentService) GetAll(
	page, limit int,
	status, docType string,
) ([]models.Document, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.documentRepo.FindAllWithFilter(offset, limit, status, docType)
}

// ============================================================
// SIGN DOCUMENT
// ============================================================
func (s *DocumentService) SignDocument(id uint, signature string) (*models.Document, error) {
	if signature == "" {
		return nil, errors.New("signature wajib diisi")
	}

	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if document.Status == "signed" {
		return nil, errors.New("document already signed")
	}

	now := time.Now()
	document.Signature = signature
	document.SignedAt = &now
	document.Status = "signed"

	if err := s.documentRepo.Update(document); err != nil {
		return nil, err
	}

	return s.documentRepo.FindByID(id)
}

// ============================================================
// SCAN QR CODE
// ============================================================
func (s *DocumentService) ScanQRCode(qrCode string) (*ScanQRResponse, error) {
	document, err := s.documentRepo.FindByQRCode(qrCode)
	if err != nil {
		document, err = s.documentRepo.FindByCode(qrCode)
		if err != nil {
			return nil, errors.New("document not found")
		}
	}

	response := &ScanQRResponse{
		Document: *document,
		User:     document.User,
	}

	if document.TransactionID != nil {
		transaction, err := s.transactionRepo.FindByID(*document.TransactionID)
		if err == nil {
			response.Transaction = transaction
		}
	}

	return response, nil
}

// ============================================================
// ✅ VALIDATE BY CODE — TERIMA SEMUA STATUS
// Cari di documents, fallback ke transactions (TANPA filter status)
// ============================================================
func (s *DocumentService) ValidateByCode(code string) (map[string]interface{}, error) {
	if code == "" {
		return nil, errors.New("kode tidak boleh kosong")
	}

	log.Printf("[ValidateByCode] === MULAI VALIDASI ===")
	log.Printf("[ValidateByCode] Code yang diterima: %q", code)

	// ============================================================
	// STEP 1: Cari di tabel documents (kode: SIG-xxx, DOC-xxx)
	// ============================================================
	var doc models.Document
	err := s.db.
		Preload("Transaction").
		Where("document_code = ?", code).
		First(&doc).Error

	if err == nil && doc.ID > 0 {
		log.Printf("[ValidateByCode] ✅ Ditemukan di DOCUMENTS: ID=%d, Status=%s",
			doc.ID, doc.Status)

		result := map[string]interface{}{
			"documentCode":  doc.DocumentCode,
			"type":          doc.Type,
			"status":        doc.Status,
			"transactionId": doc.TransactionID,
			"signedAt":      doc.SignedAt,
			"createdAt":     doc.CreatedAt,
		}

		if doc.Transaction != nil {
			result["transactionNo"] = doc.Transaction.TransactionCode
			result["itemName"] = doc.Transaction.Notes
			result["totalItems"] = doc.Transaction.TotalItems
			result["totalValue"] = doc.Transaction.TotalValue
		}

		return result, nil
	}

	log.Printf("[ValidateByCode] Tidak ada di documents, coba transactions...")

	// ============================================================
	// STEP 2: Fallback ke tabel transactions (kode: TRX-xxx)
	// ✅ TANPA filter status — semua status diterima
	// ============================================================
	var trx models.Transaction
	err = s.db.
		Preload("User").
		Where("transaction_code = ?", code).
		First(&trx).Error

	if err == nil && trx.ID > 0 {
		log.Printf("[ValidateByCode] ✅ Ditemukan di TRANSACTIONS: ID=%d, Status=%s",
			trx.ID, trx.Status)

		result := map[string]interface{}{
			"documentCode":    trx.TransactionCode,
			"transactionNo":   trx.TransactionCode,
			"type":            trx.Type,
			"status":          trx.Status, // ✅ semua status valid
			"totalItems":      trx.TotalItems,
			"totalValue":      trx.TotalValue,
			"transactionDate": trx.TransactionDate,
			"notes":           trx.Notes,
			"itemName":        trx.Notes,
		}

		if trx.User.ID > 0 {
			result["user"] = trx.User.Username
		}

		return result, nil
	}

	// ============================================================
	// STEP 3: Tidak ditemukan
	// ============================================================
	log.Printf("[ValidateByCode] ❌ TIDAK DITEMUKAN: %q", code)
	return nil, fmt.Errorf(
		"dokumen dengan kode '%s' tidak ditemukan di sistem", code,
	)
}

// ============================================================
// GENERATE QR CODE for Document
// ============================================================
func (s *DocumentService) GenerateQRCodeForDocument(documentID uint) (string, error) {
	document, err := s.documentRepo.FindByID(documentID)
	if err != nil {
		return "", err
	}

	return utils.GenerateQRCode(document.QRCode)
}

// ============================================================
// GET DOCUMENTS BY TRANSACTION
// ============================================================
func (s *DocumentService) GetDocumentsByTransaction(transactionID uint) ([]models.Document, error) {
	return s.documentRepo.FindByTransactionID(transactionID)
}