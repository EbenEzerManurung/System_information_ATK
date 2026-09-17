package controllers

import (
	"net/http"
	"strconv"
	"time"

	"atk-backend/database"
	"atk-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct{}

// GetDashboardStats - Get dashboard statistics
func (ac *AdminController) GetDashboardStats(c *gin.Context) {
	db := database.GetDB()

	// Total users
	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)

	// Total masters
	var totalMasters int64
	db.Model(&models.Master{}).Count(&totalMasters)

	// Total stock
	var totalStock int64
	db.Model(&models.Stock{}).Select("COALESCE(SUM(quantity), 0)").Scan(&totalStock)

	// Total transactions
	var totalTransactions int64
	db.Model(&models.Transaction{}).Count(&totalTransactions)

	// Low stock items
	var lowStockItems int64
	db.Model(&models.Stock{}).
		Joins("JOIN masters ON stocks.master_id = masters.id").
		Where("stocks.quantity <= masters.min_stock").
		Count(&lowStockItems)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_users":        totalUsers,
			"total_masters":      totalMasters,
			"total_stock":        totalStock,
			"total_transactions": totalTransactions,
			"low_stock_items":    lowStockItems,
		},
	})
}

// GetRecentActivities - Get recent activities
func (ac *AdminController) GetRecentActivities(c *gin.Context) {
	db := database.GetDB()

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var transactions []models.Transaction
	db.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
	})
}

// GetUserByID - Get user by ID
func (ac *AdminController) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Preload("Roles").Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// GetAllUsers - Get all users with pagination
func (ac *AdminController) GetAllUsers(c *gin.Context) {
	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	status := c.Query("status")

	offset := (page - 1) * limit

	query := db.Model(&models.User{}).Preload("Roles")

	// Search
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Get users
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"total_page": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// CreateUser - Create new user (Admin only)
func (ac *AdminController) CreateUser(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		FullName    string `json:"full_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		PhoneNumber string `json:"phone_number"`
		Status      string `json:"status"`
		RoleID      uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Check if username exists
	var existingUser models.User
	if err := db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Username already registered",
		})
		return
	}

	// Check if email exists
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to hash password",
		})
		return
	}

	status := input.Status
	if status == "" {
		status = "active"
	}

	now := time.Now()
	user := models.User{
		Username:    input.Username,
		FullName:    input.FullName,
		Email:       input.Email,
		Password:    string(hashedPassword),
		PhoneNumber: input.PhoneNumber,
		Status:      status,
		LastLogin:   &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create user: " + err.Error(),
		})
		return
	}

	// Assign role if specified
	if input.RoleID > 0 {
		var role models.Role
		if err := db.First(&role, input.RoleID).Error; err == nil {
			db.Create(&models.UserRole{
				UserID: user.ID,
				RoleID: role.ID,
			})
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser - Update user (Admin only)
func (ac *AdminController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	var input struct {
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Email != "" && input.Email != user.Email {
		// Check if email already used
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", input.Email, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Email already used by another user",
			})
			return
		}
		user.Email = input.Email
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to hash password",
			})
			return
		}
		user.Password = string(hashedPassword)
	}
	if input.PhoneNumber != "" {
		user.PhoneNumber = input.PhoneNumber
	}
	if input.Status != "" {
		user.Status = input.Status
	}

	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser - Delete user (Admin only)
func (ac *AdminController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Prevent deleting superadmin (ID 1)
	if user.ID == 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Cannot delete super admin user",
		})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// ToggleUserStatus - Toggle user status (Admin only)
func (ac *AdminController) ToggleUserStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ID",
		})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Prevent toggling superadmin
	if user.ID == 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Cannot toggle super admin user status",
		})
		return
	}

	// Toggle status
	if user.Status == "active" {
		user.Status = "inactive"
	} else {
		user.Status = "active"
	}
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update user status: " + err.Error(),
		})
		return
	}

	status := "activated"
	if user.Status != "active" {
		status = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User " + status + " successfully",
		"data":    user,
	})
}

// AssignRole - Assign role to user (Admin only)
func (ac *AdminController) AssignRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	var input struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Check if user exists
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Check if role exists
	var role models.Role
	if err := db.Where("id = ?", input.RoleID).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Role not found",
		})
		return
	}

	// Assign role
	userRole := models.UserRole{
		UserID: uint(userID),
		RoleID: input.RoleID,
	}

	if err := db.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to assign role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role assigned successfully",
	})
}

// RemoveRole - Remove role from user (Admin only)
func (ac *AdminController) RemoveRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	var input struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Check if user exists
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Remove role
	if err := db.Where("user_id = ? AND role_id = ?", userID, input.RoleID).Delete(&models.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to remove role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role removed successfully",
	})
}

// GetUserRoles - Get user roles (Admin only)
func (ac *AdminController) GetUserRoles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	db := database.GetDB()

	var roles []models.Role
	if err := db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get user roles: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// GetUserProfile - Get current user profile
func (ac *AdminController) GetUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()
	var user models.User

	if err := db.Preload("Roles").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUserProfile - Update current user profile
func (ac *AdminController) UpdateUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	var input struct {
		FullName    string `json:"full_name"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.PhoneNumber != "" {
		user.PhoneNumber = input.PhoneNumber
	}

	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// ChangePassword - Change user password
func (ac *AdminController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid old password",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to hash password",
		})
		return
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to change password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}