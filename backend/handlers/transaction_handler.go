package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"atk-backend/services"
	"atk-backend/utils"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		transactionService: services.NewTransactionService(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (h *TransactionHandler) Create(c *gin.Context) {
	var req services.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	userID := c.GetUint("user_id")
	transaction, err := h.transactionService.Create(req, userID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "Transaction created successfully", transaction)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	transaction, err := h.transactionService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Transaction not found")
		return
	}

	utils.SuccessResponse(c, "Transaction found", transaction)
}

// ============================================================
// GET ALL
// ============================================================
func (h *TransactionHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sizeStr := c.DefaultQuery("size", c.DefaultQuery("limit", "10"))
	size, _ := strconv.Atoi(sizeStr)
	search := c.Query("search")
	status := c.Query("status")
	transType := c.Query("type")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	transactions, total, err := h.transactionService.GetAllWithSearch(
		page, size, search, status, transType, uint(userID),
	)
	if err != nil {
		log.Printf("[TransactionGetAll] Error: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch transactions")
		return
	}

	totalPages := (int(total) + size - 1) / size
	if totalPages < 1 {
		totalPages = 1
	}

	utils.SuccessResponse(c, "Transactions fetched successfully", gin.H{
		"data": transactions,
		"meta": gin.H{
			"page":       page,
			"size":       size,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// ============================================================
// UPDATE
// ============================================================
func (h *TransactionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	transaction, err := h.transactionService.Update(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaction updated successfully", transaction)
}

// ============================================================
// SUBMIT
// ============================================================
func (h *TransactionHandler) Submit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.transactionService.Submit(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaction submitted successfully", nil)
}

// ============================================================
// APPROVE — ✅ TERIMA SIGNATURE DARI BODY
// ============================================================
func (h *TransactionHandler) Approve(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	// ✅ Bind request body untuk dapat signature & notes
	var req services.ApproveTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Boleh kosong, tidak wajib body
		log.Printf("[TransactionApprove] Bind warning (oke jika kosong): %v", err)
		req = services.ApproveTransactionRequest{}
	}

	approverID := c.GetUint("user_id")

	log.Printf("[TransactionApprove] ID=%d, ApproverID=%d, HasSignature=%v",
		id, approverID, req.Signature != "" || req.SignatureData != "")

	transaction, err := h.transactionService.Approve(uint(id), approverID, req)
	if err != nil {
		log.Printf("[TransactionApprove] Error: %v", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaksi berhasil disetujui", transaction)
}

// ============================================================
// REJECT — ✅ METHOD BARU
// ============================================================
func (h *TransactionHandler) Reject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.RejectTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransactionReject] Bind warning: %v", err)
		req = services.RejectTransactionRequest{}
	}

	approverID := c.GetUint("user_id")

	log.Printf("[TransactionReject] ID=%d, ApproverID=%d, Reason=%q",
		id, approverID, req.Reason)

	transaction, err := h.transactionService.Reject(uint(id), approverID, req)
	if err != nil {
		log.Printf("[TransactionReject] Error: %v", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaksi ditolak", transaction)
}

// ============================================================
// COMPLETE
// ============================================================
func (h *TransactionHandler) Complete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.transactionService.Complete(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaction completed successfully", nil)
}

// ============================================================
// CANCEL
// ============================================================
func (h *TransactionHandler) Cancel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.transactionService.Cancel(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaction cancelled successfully", nil)
}

// ============================================================
// DELETE
// ============================================================
func (h *TransactionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.transactionService.Delete(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Transaction deleted successfully", nil)
}

// ============================================================
// EXPORT EXCEL
// ============================================================
func (h *TransactionHandler) ExportExcel(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	transType := c.Query("type")

	log.Printf("[ExportTransaction] Request: search='%s', status='%s', type='%s'",
		search, status, transType)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ExportTransaction] PANIC RECOVERED: %v", r)
			utils.InternalServerErrorResponse(c, "Internal error saat export")
		}
	}()

	buf, err := h.transactionService.ExportToExcel(search, status, transType)
	if err != nil {
		log.Printf("[ExportTransaction] Service error: %v", err)
		utils.InternalServerErrorResponse(c, "Gagal export: "+err.Error())
		return
	}

	if buf == nil || buf.Len() == 0 {
		log.Printf("[ExportTransaction] Buffer kosong")
		utils.InternalServerErrorResponse(c, "Tidak ada data untuk di-export")
		return
	}

	bufLen := buf.Len()
	log.Printf("[ExportTransaction] Buffer size: %d bytes", bufLen)

	// Verifikasi magic bytes XLSX (PK)
	bytes := buf.Bytes()
	if len(bytes) < 2 || bytes[0] != 0x50 || bytes[1] != 0x4B {
		log.Printf("[ExportTransaction] Bukan XLSX valid")
		utils.InternalServerErrorResponse(c, "File bukan format XLSX")
		return
	}

	filename := "transactions_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(bufLen))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")

	c.Data(http.StatusOK, contentType, bytes)
	log.Printf("[ExportTransaction] Selesai, terkirim %d bytes", bufLen)
}

// ============================================================
// IMPORT EXCEL
// ============================================================
func (h *TransactionHandler) ImportExcel(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File tidak ditemukan", err.Error())
		return
	}

	log.Printf("[ImportTransaction] File: %s, size: %d bytes",
		fileHeader.Filename, fileHeader.Size)

	userID := c.GetUint("user_id")

	result, err := h.transactionService.ImportFromExcel(fileHeader, userID)
	if err != nil {
		log.Printf("[ImportTransaction] Error: %v", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	log.Printf("[ImportTransaction] Done: imported=%d, failed=%d",
		result.Imported, result.Failed)

	utils.SuccessResponse(c, "Import berhasil", gin.H{
		"imported": result.Imported,
		"failed":   result.Failed,
		"errors":   result.Errors,
	})
}