package handlers

import (
	"log"
	"strconv"

	"atk-backend/services"
	"atk-backend/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService *services.DocumentService
}

func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{
		documentService: services.NewDocumentService(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (h *DocumentHandler) Create(c *gin.Context) {
	var req services.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	// Kalau UserID kosong, isi dari token
	if req.UserID == 0 {
		req.UserID = c.GetUint("user_id")
	}

	document, err := h.documentService.Create(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "Document created successfully", document)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *DocumentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	document, err := h.documentService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Document not found")
		return
	}

	utils.SuccessResponse(c, "Document found", document)
}

// ============================================================
// GET BY CODE
// ============================================================
func (h *DocumentHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		utils.BadRequestResponse(c, "Invalid code", nil)
		return
	}

	document, err := h.documentService.GetByCode(code)
	if err != nil {
		utils.NotFoundResponse(c, "Document not found")
		return
	}

	utils.SuccessResponse(c, "Document found", document)
}

// ============================================================
// GET ALL — pagination + filter
// Response konsisten: { success, message, data: { data, meta } }
// ============================================================
func (h *DocumentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sizeStr := c.DefaultQuery("size", c.DefaultQuery("limit", "10"))
	size, _ := strconv.Atoi(sizeStr)
	status := c.Query("status")
	docType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	documents, total, err := h.documentService.GetAll(page, size, status, docType)
	if err != nil {
		log.Printf("[DocumentGetAll] Error: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch documents")
		return
	}

	totalPages := (int(total) + size - 1) / size
	if totalPages < 1 {
		totalPages = 1
	}

	utils.SuccessResponse(c, "Documents fetched successfully", gin.H{
		"data": documents,
		"meta": gin.H{
			"page":       page,
			"size":       size,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// ============================================================
// SIGN DOCUMENT — tanda tangan langsung ke document
// ============================================================
func (h *DocumentHandler) SignDocument(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req struct {
		Signature     string `json:"signature"`
		SignatureData string `json:"signatureData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	// Terima signature dari salah satu field
	sig := req.Signature
	if sig == "" {
		sig = req.SignatureData
	}
	if sig == "" {
		utils.BadRequestResponse(c, "Signature wajib diisi", nil)
		return
	}

	document, err := h.documentService.SignDocument(uint(id), sig)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Document signed successfully", document)
}

// ============================================================
// SCAN QR CODE
// ============================================================
func (h *DocumentHandler) ScanQRCode(c *gin.Context) {
	var req struct {
		QRCode string `json:"qrCode"`
		Code   string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	// Terima qrCode atau code
	qr := req.QRCode
	if qr == "" {
		qr = req.Code
	}
	// Bisa juga dari query param
	if qr == "" {
		qr = c.Query("code")
	}
	if qr == "" {
		utils.BadRequestResponse(c, "QR code atau document code wajib diisi", nil)
		return
	}

	result, err := h.documentService.ScanQRCode(qr)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "QR Code scanned successfully", result)
}

// ============================================================
// ✅ VALIDATE DOCUMENT — untuk scan QR
// Route: GET /api/documents/validate?code=XXX
// Cari di tabel documents ATAU transactions
// ============================================================
func (h *DocumentHandler) ValidateDocument(c *gin.Context) {
	// Terima dari query param "code" atau "documentCode"
	code := c.Query("code")
	if code == "" {
		code = c.Query("documentCode")
	}
	if code == "" {
		utils.BadRequestResponse(c, "Parameter 'code' wajib diisi", nil)
		return
	}

	log.Printf("[ValidateDocument] Validating code: %s", code)

	// Panggil service untuk validasi
	result, err := h.documentService.ValidateByCode(code)
	if err != nil {
		log.Printf("[ValidateDocument] Not found: %v", err)
		// Return success=true dengan valid=false (bukan error HTTP)
		// supaya frontend bisa tampilkan pesan
		utils.SuccessResponse(c, "Dokumen tidak ditemukan", gin.H{
			"valid":   false,
			"message": err.Error(),
		})
		return
	}

	log.Printf("[ValidateDocument] Valid: %v", result["documentCode"])

	utils.SuccessResponse(c, "Dokumen valid", gin.H{
		"valid":    true,
		"document": result,
	})
}

// ============================================================
// GENERATE QR CODE untuk Document
// ============================================================
func (h *DocumentHandler) GenerateQRCode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	qrImage, err := h.documentService.GenerateQRCodeForDocument(uint(id))
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "QR Code generated successfully", gin.H{
		"qrImage":  qrImage,
		"qr_image": qrImage, // untuk backward compat
	})
}