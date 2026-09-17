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

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (h *RoleHandler) Create(c *gin.Context) {
	var req services.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	role, err := h.roleService.Create(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "Role created successfully", role)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Role not found")
		return
	}

	utils.SuccessResponse(c, "Role found", role)
}

// ============================================================
// GET ALL
// ============================================================
func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch roles")
		return
	}

	utils.SuccessResponse(c, "Roles fetched successfully", roles)
}

// ============================================================
// UPDATE
// ============================================================
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	role, err := h.roleService.Update(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Role updated successfully", role)
}

// ============================================================
// DELETE
// ============================================================
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	err = h.roleService.Delete(uint(id))
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Role deleted successfully", nil)
}

// ============================================================
// IMPORT EXCEL
// Format kolom (baris 1 = header):
// A: Nama Role | B: Deskripsi
// ============================================================
func (h *RoleHandler) ImportExcel(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File tidak ditemukan", err.Error())
		return
	}

	log.Printf("[ImportRole] File: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

	result, err := h.roleService.ImportFromExcel(fileHeader)
	if err != nil {
		log.Printf("[ImportRole] ERROR: %v", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	log.Printf("[ImportRole] Done: imported=%d, failed=%d", result.Imported, result.Failed)

	utils.SuccessResponse(c, "Import berhasil", gin.H{
		"imported": result.Imported,
		"failed":   result.Failed,
		"errors":   result.Errors,
	})
}

// ============================================================
// EXPORT EXCEL
// ============================================================
func (h *RoleHandler) ExportExcel(c *gin.Context) {
	search := c.Query("search")

	log.Printf("[ExportRole] Request search='%s'", search)

	buf, err := h.roleService.ExportToExcel(search)
	if err != nil {
		log.Printf("[ExportRole] ERROR: %v", err)
		utils.InternalServerErrorResponse(c, "Gagal export data: "+err.Error())
		return
	}

	if buf == nil {
		log.Printf("[ExportRole] Buffer nil")
		utils.InternalServerErrorResponse(c, "Buffer kosong")
		return
	}

	bufLen := buf.Len()
	log.Printf("[ExportRole] Buffer size: %d bytes", bufLen)

	if bufLen == 0 {
		log.Printf("[ExportRole] Buffer kosong (0 bytes)")
		utils.InternalServerErrorResponse(c, "Tidak ada data untuk di-export")
		return
	}

	bytes := buf.Bytes()
	if len(bytes) < 4 || bytes[0] != 0x50 || bytes[1] != 0x4B {
		log.Printf("[ExportRole] Bukan XLSX valid. First bytes: %x", bytes[:4])
		utils.InternalServerErrorResponse(c, "File bukan format XLSX")
		return
	}

	filename := "roles_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(bufLen))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")

	c.Data(http.StatusOK, contentType, bytes)
}