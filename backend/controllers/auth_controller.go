package controllers

import (
	"net/http"
	"time"

	"atk-backend/database"
	"atk-backend/models"
	"atk-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (ac *AuthController) Register(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		FullName    string `json:"full_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	// Check if username already exists
	var existingUser models.User
	if err := db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already registered"})
		return
	}

	// Check if email already exists
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	now := time.Now()
	user := models.User{
		Username:    input.Username,
		FullName:    input.FullName,
		Email:       input.Email,
		Password:    string(hashedPassword),
		PhoneNumber: input.PhoneNumber,
		Status:      "active",
		LastLogin:   &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"status":       user.Status,
		},
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// Try to find by email if username not found
		if err := db.Where("email = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Check if user is active
	if user.Status != "active" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not active"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	now := time.Now()
	db.Model(&user).Update("last_login", &now)

	// Get user roles
	var roles []models.Role
	db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", user.ID).
		Find(&roles)
	user.Roles = roles

	// Get primary role
	role := ""
	if len(roles) > 0 {
		role = roles[0].Name
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"status":       user.Status,
			"roles":        roles,
		},
	})
}