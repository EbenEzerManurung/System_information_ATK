package handlers

import (
	"strconv"
    "time"    
	"atk-backend/services"
	"atk-backend/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// ============================================================
// CRUD USER
// ============================================================

// Create - POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.Create(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "User created successfully", user)
}

// GetByID - GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User found", user)
}

// GetByUsername - GET /api/users/username/:username
func (h *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetByUsername(username)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User found", user)
}

// GetAll - GET /api/users?page=1&limit=10&search=xxx
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	users, total, err := h.userService.GetAll(page, limit, search)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch users")
		return
	}

	utils.SuccessWithMeta(c, "Users fetched successfully", users, gin.H{
		"page":      page,
		"limit":     limit,
		"total":     total,
		"totalPage": (total + int64(limit) - 1) / int64(limit),
	})
}

// Update - PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.Update(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "User updated successfully", user)
}

// Delete - DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID", nil)
		return
	}

	currentUserID := c.GetUint("user_id")
	if uint(id) == currentUserID {
		utils.BadRequestResponse(c, "Cannot delete your own account", nil)
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}

// ============================================================
// PROFILE (untuk user yang sedang login)
// ============================================================

// GetProfile - GET /api/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile found", user)
}

// UpdateProfile - PUT /api/profile
// Body: { "fullName": "...", "email": "..." }
// Catatan: frontend Angular mengirim camelCase. Jika service Anda
// menggunakan snake_case, tambahkan binding tag di UpdateProfileRequest.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req services.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user)
}

// ChangePassword - POST /api/profile/change-password
// Body: { "currentPassword": "...", "newPassword": "...", "confirmPassword": "..." }
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req services.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, req); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Password changed successfully", nil)
}

// ============================================================
// ROLE ASSIGNMENT
// ============================================================

// AssignRole - POST /api/users/:userId/roles
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID", nil)
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	if err := h.userService.AssignRole(uint(userID), req.RoleID); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Role assigned successfully", nil)
}

// RemoveRole - DELETE /api/users/:userId/roles
func (h *UserHandler) RemoveRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID", nil)
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	if err := h.userService.RemoveRole(uint(userID), req.RoleID); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Role removed successfully", nil)
}

// GetUserRoles - GET /api/users/:userId/roles
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID", nil)
		return
	}

	roles, err := h.userService.GetUserRoles(uint(userID))
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "User roles fetched successfully", roles)
}

// ============================================================
// IMPORT / EXPORT EXCEL (untuk frontend Angular)
// ============================================================

// ImportExcel - POST /api/users/import
// Menerima multipart/form-data dengan field "file"
func (h *UserHandler) ImportExcel(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File tidak ditemukan", err.Error())
		return
	}

	result, err := h.userService.ImportFromExcel(fileHeader)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Import berhasil", gin.H{
		"imported": result.Imported,
		"failed":   result.Failed,
		"errors":   result.Errors,
	})
}

// ExportExcel - GET /api/users/export?search=xxx
// Mengembalikan file .xlsx sebagai attachment
func (h *UserHandler) ExportExcel(c *gin.Context) {
	search := c.Query("search")

	buf, err := h.userService.ExportToExcel(search)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Gagal export data")
		return
	}

	filename := "users_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Data(200, c.Writer.Header().Get("Content-Type"), buf.Bytes())
}