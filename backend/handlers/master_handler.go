package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"atk-backend/services"
	"atk-backend/utils"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	masterService *services.MasterService
	excelService  *services.ExcelService
}

func NewMasterHandler() *MasterHandler {
	return &MasterHandler{
		masterService: services.NewMasterService(),
		excelService:  services.NewExcelService(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (h *MasterHandler) Create(c *gin.Context) {
	var req services.CreateMasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	master, err := h.masterService.Create(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "Master created successfully", master)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *MasterHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	master, err := h.masterService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Master not found")
		return
	}

	utils.SuccessResponse(c, "Master found", master)
}

// ============================================================
// GET ALL — dengan search + filter category + pagination
// Format response konsisten dengan Stock/Role/User
// ============================================================
func (h *MasterHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// Terima "size" ATAU "limit"
	sizeStr := c.DefaultQuery("size", c.DefaultQuery("limit", "10"))
	size, _ := strconv.Atoi(sizeStr)
	search := c.Query("search")
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	masters, total, err := h.masterService.GetAll(page, size, search, category)
	if err != nil {
		log.Printf("[MasterGetAll] Error: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch masters")
		return
	}

	totalPages := (int(total) + size - 1) / size
	if totalPages < 1 {
		totalPages = 1
	}

	// Format konsisten: { data: { data: [...], meta: {...} } }
	utils.SuccessResponse(c, "Masters fetched successfully", gin.H{
		"data": masters,
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
func (h *MasterHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.UpdateMasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error(), nil)
		return
	}

	master, err := h.masterService.Update(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Master updated successfully", master)
}

// ============================================================
// DELETE
// ============================================================
func (h *MasterHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	if err := h.masterService.Delete(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Master deleted successfully", nil)
}

// ============================================================
// EXPORT EXCEL — dengan filter search + category
// ============================================================
func (h *MasterHandler) ExportExcel(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")

	log.Printf("[ExportMaster] Request: search='%s', category='%s'", search, category)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ExportMaster] PANIC RECOVERED: %v", r)
			utils.InternalServerErrorResponse(c, "Internal error saat export")
		}
	}()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	// Panggil excel service — sesuaikan signature kalau perlu filter
	buf, _, err := h.excelService.ExportMasterData(ctx)
	if err != nil {
		log.Printf("[ExportMaster] Service error: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to export data: "+err.Error())
		return
	}

	if buf == nil || buf.Len() == 0 {
		log.Printf("[ExportMaster] Buffer kosong")
		utils.InternalServerErrorResponse(c, "Tidak ada data untuk di-export")
		return
	}

	bufLen := buf.Len()
	log.Printf("[ExportMaster] Buffer size: %d bytes", bufLen)

	// Verifikasi magic bytes XLSX (PK)
	bytes := buf.Bytes()
	if len(bytes) < 2 || bytes[0] != 0x50 || bytes[1] != 0x4B {
		log.Printf("[ExportMaster] Bukan XLSX valid")
		utils.InternalServerErrorResponse(c, "File bukan format XLSX")
		return
	}

	filename := "master_data_" + time.Now().Format("20060102_150405") + ".xlsx"
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(bufLen))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")

	c.Data(http.StatusOK, contentType, bytes)
	log.Printf("[ExportMaster] Selesai, terkirim %d bytes", bufLen)
}

// ============================================================
// IMPORT EXCEL
// ============================================================
func (h *MasterHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File is required", err.Error())
		return
	}

	log.Printf("[ImportMaster] File: %s, size: %d bytes", file.Filename, file.Size)

	if file.Size > 50*1024*1024 {
		utils.BadRequestResponse(c, "File too large (max 50MB)", nil)
		return
	}

	src, err := file.Open()
	if err != nil {
		utils.BadRequestResponse(c, "Failed to open file", err.Error())
		return
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		utils.BadRequestResponse(c, "Failed to read file", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	imported, errors, err := h.excelService.ImportMasterData(ctx, buf)
	if err != nil {
		log.Printf("[ImportMaster] Error: %v", err)
		utils.BadRequestResponse(c, "Import failed: "+err.Error(), nil)
		return
	}

	log.Printf("[ImportMaster] Done: imported=%d, errors=%d", imported, len(errors))

	utils.SuccessResponse(c,
		"Successfully imported data",
		gin.H{
			"imported": imported,
			"failed":   len(errors),
			"errors":   errors,
		})
}

// ============================================================
// DOWNLOAD TEMPLATE
// ============================================================
func (h *MasterHandler) DownloadTemplate(c *gin.Context) {
	buf, err := h.excelService.ExportTemplate()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate template")
		return
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=import_template.xlsx")
	c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")

	c.Data(http.StatusOK, contentType, buf.Bytes())
}