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

// ============================================================
// STOCK HANDLER
// ============================================================
type StockHandler struct {
	stockService *services.StockService
}

func NewStockHandler() *StockHandler {
	return &StockHandler{
		stockService: services.NewStockService(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (h *StockHandler) Create(c *gin.Context) {
	var req services.CreateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	stock, err := h.stockService.Create(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "Stock created successfully", stock)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *StockHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	stock, err := h.stockService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Stock not found")
		return
	}

	utils.SuccessResponse(c, "Stock found", stock)
}

// ============================================================
// GET BY MASTER ID
// ============================================================
func (h *StockHandler) GetByMasterID(c *gin.Context) {
	masterID, err := strconv.Atoi(c.Param("masterId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid master ID", nil)
		return
	}

	stock, err := h.stockService.GetByMasterID(uint(masterID))
	if err != nil {
		utils.NotFoundResponse(c, "Stock not found")
		return
	}

	utils.SuccessResponse(c, "Stock found", stock)
}

// ============================================================
// GET ALL — dengan search + filter status + Preload Master
// ============================================================
func (h *StockHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sizeStr := c.DefaultQuery("size", c.DefaultQuery("limit", "10"))
	size, _ := strconv.Atoi(sizeStr)
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	stocks, total, err := h.stockService.GetAllWithMaster(page, size, search, status)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch stocks: "+err.Error())
		return
	}

	totalPages := (int(total) + size - 1) / size
	if totalPages < 1 {
		totalPages = 1
	}

	utils.SuccessResponse(c, "Stocks fetched successfully", gin.H{
		"data": stocks,
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
func (h *StockHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	stock, err := h.stockService.Update(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Stock updated successfully", stock)
}

// ============================================================
// DELETE
// ============================================================
func (h *StockHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.stockService.Delete(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Stock deleted successfully", nil)
}

// ============================================================
// GET LOW STOCK
// ============================================================
func (h *StockHandler) GetLowStock(c *gin.Context) {
	minStock, _ := strconv.Atoi(c.DefaultQuery("min_stock", "10"))

	stocks, err := h.stockService.GetLowStock(minStock)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch low stock items")
		return
	}

	utils.SuccessResponse(c, "Low stock items fetched successfully", stocks)
}

// ============================================================
// GET STOCK SUMMARY
// ============================================================
func (h *StockHandler) GetStockSummary(c *gin.Context) {
	summary, err := h.stockService.GetStockSummary()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch stock summary")
		return
	}

	utils.SuccessResponse(c, "Stock summary fetched successfully", summary)
}

// ============================================================
// ADJUST STOCK
// ============================================================
func (h *StockHandler) AdjustStock(c *gin.Context) {
	masterID, err := strconv.Atoi(c.Param("masterId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid master ID", nil)
		return
	}

	var req struct {
		Quantity int    `json:"quantity" binding:"required"`
		Notes    string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	stock, err := h.stockService.AdjustStock(uint(masterID), req.Quantity, req.Notes)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Stock adjusted successfully", stock)
}

// ============================================================
// IMPORT EXCEL
// ============================================================
func (h *StockHandler) ImportExcel(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File tidak ditemukan", err.Error())
		return
	}

	log.Printf("[ImportExcel] Received file: %s, size: %d bytes",
		fileHeader.Filename, fileHeader.Size)

	result, err := h.stockService.ImportFromExcel(fileHeader)
	if err != nil {
		log.Printf("[ImportExcel] ERROR: %v", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	log.Printf("[ImportExcel] Success: imported=%d, updated=%d, failed=%d",
		result.Imported, result.Updated, result.Failed)

	utils.SuccessResponse(c, "Import berhasil", gin.H{
		"imported": result.Imported,
		"updated":  result.Updated,
		"failed":   result.Failed,
		"errors":   result.Errors,
	})
}

// ============================================================
// EXPORT EXCEL — dengan anti-cache + validasi buffer
// ============================================================
func (h *StockHandler) ExportExcel(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	log.Printf("[ExportExcel] Request: search='%s', status='%s'", search, status)

	// Generate Excel buffer
	buf, err := h.stockService.ExportToExcel(search, status)
	if err != nil {
		log.Printf("[ExportExcel] ERROR: %v", err)
		utils.InternalServerErrorResponse(c, "Gagal export data: "+err.Error())
		return
	}

	// Validasi buffer
	if buf == nil {
		log.Printf("[ExportExcel] ERROR: buffer nil")
		utils.InternalServerErrorResponse(c, "Buffer kosong")
		return
	}

	bufLen := buf.Len()
	log.Printf("[ExportExcel] Buffer size: %d bytes", bufLen)

	if bufLen == 0 {
		log.Printf("[ExportExcel] WARNING: buffer kosong (0 bytes)")
		utils.InternalServerErrorResponse(c, "Tidak ada data untuk di-export")
		return
	}

	// Verifikasi magic bytes XLSX (PK\x03\x04)
	bytes := buf.Bytes()
	if len(bytes) < 4 {
		log.Printf("[ExportExcel] ERROR: buffer terlalu pendek")
		utils.InternalServerErrorResponse(c, "File tidak valid")
		return
	}
	if bytes[0] != 0x50 || bytes[1] != 0x4B {
		log.Printf("[ExportExcel] ERROR: bukan file ZIP/XLSX. First bytes: %x", bytes[:4])
		utils.InternalServerErrorResponse(c, "File bukan format XLSX")
		return
	}

	// Header untuk file download
	filename := "stock_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(bufLen))

	// 🔑 Anti-cache headers — WAJIB untuk file download
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 🔑 Expose header agar frontend bisa baca filename
	c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")

	c.Data(http.StatusOK, contentType, bytes)
}