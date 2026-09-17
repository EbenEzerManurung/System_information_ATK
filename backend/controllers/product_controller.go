package controllers

import (
	"net/http"
	"strconv"

	"atk-backend/database"
	"atk-backend/models"

	"github.com/gin-gonic/gin"
)

type ProductResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// GetProductList - Get all masters with filters (renamed from GetProducts)
func GetProductList(c *gin.Context) {
	var masters []models.Master
	var total int64

	db := database.GetDB()
	query := db.Model(&models.Master{}).Preload("Stocks")

	// Filter by category
	category := c.Query("category")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Search by name or code
	search := c.Query("search")
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	query.Count(&total)

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	query = query.Offset(offset).Limit(limit)
	query = query.Order("created_at DESC")

	if err := query.Find(&masters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to get masters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Masters retrieved successfully",
		Data:    masters,
		Total:   total,
	})
}

// GetProductByID - Get single master by ID (renamed from GetProduct)
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var master models.Master

	db := database.GetDB()
	if err := db.Preload("Stocks").First(&master, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Master not found",
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Master retrieved successfully",
		Data:    master,
	})
}

// CreateNewProduct - Create new master (renamed from CreateProduct)
func CreateNewProduct(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "super_admin" && role != "admin" {
		c.JSON(http.StatusForbidden, ProductResponse{
			Success: false,
			Message: "Only admin can create masters",
		})
		return
	}

	var input struct {
		Code         string  `json:"code" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		Category     string  `json:"category" binding:"required"`
		Description  string  `json:"description"`
		Unit         string  `json:"unit" binding:"required"`
		MinStock     int     `json:"min_stock"`
		MaxStock     int     `json:"max_stock"`
		Price        float64 `json:"price" binding:"required,min=0"`
		InitialStock int     `json:"initial_stock"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	var existingMaster models.Master
	if err := db.Where("code = ?", input.Code).First(&existingMaster).Error; err == nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Master code already exists",
		})
		return
	}

	master := models.Master{
		Code:        input.Code,
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
		Unit:        input.Unit,
		MinStock:    input.MinStock,
		MaxStock:    input.MaxStock,
		Price:       input.Price,
	}

	if err := db.Create(&master).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to create master: " + err.Error(),
		})
		return
	}

	if input.InitialStock > 0 {
		stock := models.Stock{
			MasterID: master.ID,
			Quantity: input.InitialStock,
			Location: "Gudang Utama",
			Status:   "available",
		}
		db.Create(&stock)
	}

	c.JSON(http.StatusCreated, ProductResponse{
		Success: true,
		Message: "Master created successfully",
		Data:    master,
	})
}

// UpdateProductByID - Update existing master (renamed from UpdateProduct)
func UpdateProductByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" {
		c.JSON(http.StatusForbidden, ProductResponse{
			Success: false,
			Message: "Only admin can update masters",
		})
		return
	}

	db := database.GetDB()
	var master models.Master
	if err := db.First(&master, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Master not found",
		})
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Unit        string  `json:"unit"`
		MinStock    int     `json:"min_stock"`
		MaxStock    int     `json:"max_stock"`
		Price       float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	if input.Name != "" {
		master.Name = input.Name
	}
	if input.Category != "" {
		master.Category = input.Category
	}
	if input.Description != "" {
		master.Description = input.Description
	}
	if input.Unit != "" {
		master.Unit = input.Unit
	}
	if input.MinStock >= 0 {
		master.MinStock = input.MinStock
	}
	if input.MaxStock >= 0 {
		master.MaxStock = input.MaxStock
	}
	if input.Price >= 0 {
		master.Price = input.Price
	}

	if err := db.Save(&master).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to update master: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Master updated successfully",
		Data:    master,
	})
}

// DeleteProductByID - Delete master (renamed from DeleteProduct)
func DeleteProductByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" {
		c.JSON(http.StatusForbidden, ProductResponse{
			Success: false,
			Message: "Only admin can delete masters",
		})
		return
	}

	db := database.GetDB()
	var master models.Master
	if err := db.First(&master, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Master not found",
		})
		return
	}

	if err := db.Delete(&master).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to delete master: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Master deleted successfully",
	})
}

// ToggleProductStatusByID - Toggle master stock status (renamed from ToggleProductStatus)
func ToggleProductStatusByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, ProductResponse{
			Success: false,
			Message: "Only admin or manager can toggle master status",
		})
		return
	}

	db := database.GetDB()
	var master models.Master
	if err := db.First(&master, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Master not found",
		})
		return
	}

	var stock models.Stock
	if err := db.Where("master_id = ?", master.ID).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Stock not found for this master",
		})
		return
	}

	if stock.Status == "available" {
		stock.Status = "unavailable"
	} else {
		stock.Status = "available"
	}

	if err := db.Save(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to toggle master status: " + err.Error(),
		})
		return
	}

	status := "available"
	if stock.Status != "available" {
		status = "unavailable"
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Master status changed to " + status,
		Data:    master,
	})
}

// GetProductStockByID - Get stock information for a master (renamed from GetProductStock)
func GetProductStockByID(c *gin.Context) {
	id := c.Param("id")
	var stock models.Stock

	db := database.GetDB()
	if err := db.Preload("Master").Where("master_id = ?", id).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Stock not found for this master",
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Stock retrieved successfully",
		Data:    stock,
	})
}

// UpdateProductStockByID - Update stock for a master (renamed from UpdateProductStock)
func UpdateProductStockByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, ProductResponse{
			Success: false,
			Message: "Only admin or manager can update stock",
		})
		return
	}

	var input struct {
		Quantity int    `json:"quantity" binding:"required"`
		Location string `json:"location"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	var stock models.Stock
	if err := db.Where("master_id = ?", id).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Stock not found",
		})
		return
	}

	if input.Quantity >= 0 {
		stock.Quantity = input.Quantity
	}
	if input.Location != "" {
		stock.Location = input.Location
	}
	if input.Status != "" {
		stock.Status = input.Status
	}

	if err := db.Save(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProductResponse{
			Success: false,
			Message: "Failed to update stock: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Stock updated successfully",
		Data:    stock,
	})
}