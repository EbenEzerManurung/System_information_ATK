package controllers

import (
	"net/http"
	"strconv"
	"time"

	"atk-backend/database"
	"atk-backend/models"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// CreateTransaction - Create new transaction
func CreateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	var input struct {
		Type  string `json:"type" binding:"required,oneof=IN OUT RETURN ADJUSTMENT"`
		Notes string `json:"notes"`
		Items []struct {
			MasterID uint `json:"master_id" binding:"required"`
			Quantity int  `json:"quantity" binding:"required,gt=0"`
			Notes    string `json:"notes"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Check permission for OUT transactions (only admin/manager can do OUT)
	if input.Type == "OUT" && role != "super_admin" && role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "Only admin or manager can create OUT transactions",
		})
		return
	}

	// Generate transaction code
	code := "TRX-" + time.Now().Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%10000, 10)

	transaction := models.Transaction{
		TransactionCode: code,
		Type:            input.Type,
		Status:          "draft",
		UserID:          userID,
		Notes:           input.Notes,
		TransactionDate: time.Now(),
	}

	var totalItems int
	var totalValue float64
	var items []models.TransactionItem

	// Process each item
	for _, item := range input.Items {
		var master models.Master
		if err := db.First(&master, item.MasterID).Error; err != nil {
			c.JSON(http.StatusNotFound, TransactionResponse{
				Success: false,
				Message: "Master item not found",
			})
			return
		}

		// Check stock for OUT transactions
		if input.Type == "OUT" {
			var stock models.Stock
			if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err != nil {
				c.JSON(http.StatusBadRequest, TransactionResponse{
					Success: false,
					Message: "Stock not found for item: " + master.Name,
				})
				return
			}
			if stock.Quantity < item.Quantity {
				c.JSON(http.StatusBadRequest, TransactionResponse{
					Success: false,
					Message: "Insufficient stock for item: " + master.Name + " (available: " + strconv.Itoa(stock.Quantity) + ")",
				})
				return
			}
		}

		subTotal := float64(item.Quantity) * master.Price
		totalItems += item.Quantity
		totalValue += subTotal

		transactionItem := models.TransactionItem{
			MasterID: item.MasterID,
			Quantity: item.Quantity,
			Price:    master.Price,
			SubTotal: subTotal,
			Notes:    item.Notes,
		}
		items = append(items, transactionItem)
	}

	transaction.TotalItems = totalItems
	transaction.TotalValue = totalValue
	transaction.Items = items

	// Save transaction
	if err := db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to create transaction: " + err.Error(),
		})
		return
	}

	// Update stock for OUT transactions
	if input.Type == "OUT" {
		for _, item := range items {
			var stock models.Stock
			if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err == nil {
				stock.Quantity -= item.Quantity
				stock.LastUpdated = time.Now()
				db.Save(&stock)
			}
		}
	}

	c.JSON(http.StatusCreated, TransactionResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    transaction,
	})
}

// GetTransactions - Get all transactions
func GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	transType := c.Query("type")

	offset := (page - 1) * limit

	query := db.Model(&models.Transaction{}).
		Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if transType != "" {
		query = query.Where("type = ?", transType)
	}
	if role != "super_admin" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var transactions []models.Transaction
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to get transactions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    transactions,
		Total:   total,
	})
}

// GetTransaction - Get transaction by ID
func GetTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master").
		First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && transaction.UserID != userID {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "You don't have permission to view this transaction",
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}

// UpdateTransaction - Update transaction
func UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && transaction.UserID != userID {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "You don't have permission to update this transaction",
		})
		return
	}

	if transaction.Status != "draft" {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Only draft transactions can be updated",
		})
		return
	}

	var input struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	if input.Notes != "" {
		transaction.Notes = input.Notes
	}

	if err := db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to update transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction updated successfully",
		Data:    transaction,
	})
}

// SubmitTransaction - Submit transaction for approval
func SubmitTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && transaction.UserID != userID {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "You don't have permission to submit this transaction",
		})
		return
	}

	if transaction.Status != "draft" {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Only draft transactions can be submitted",
		})
		return
	}

	transaction.Status = "pending"
	if err := db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to submit transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction submitted successfully",
		Data:    transaction,
	})
}

// ApproveTransaction - Approve transaction
func ApproveTransaction(c *gin.Context) {
	id := c.Param("id")
	approverID := c.GetUint("user_id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "Only admin or manager can approve transactions",
		})
		return
	}

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	if transaction.Status != "pending" {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Only pending transactions can be approved",
		})
		return
	}

	// Update stock for OUT transactions when approved
	if transaction.Type == "OUT" {
		var items []models.TransactionItem
		db.Where("transaction_id = ?", transaction.ID).Find(&items)

		for _, item := range items {
			var stock models.Stock
			if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err == nil {
				if stock.Quantity < item.Quantity {
					c.JSON(http.StatusBadRequest, TransactionResponse{
						Success: false,
						Message: "Insufficient stock for item ID: " + strconv.Itoa(int(item.MasterID)),
					})
					return
				}
				stock.Quantity -= item.Quantity
				stock.LastUpdated = time.Now()
				db.Save(&stock)
			}
		}
	}

	transaction.Status = "approved"
	transaction.ApprovedBy = &approverID
	now := time.Now()
	transaction.ApprovedAt = &now

	if err := db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to approve transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction approved successfully",
		Data:    transaction,
	})
}

// CompleteTransaction - Complete transaction
func CompleteTransaction(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "Only admin can complete transactions",
		})
		return
	}

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	if transaction.Status != "approved" {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Only approved transactions can be completed",
		})
		return
	}

	transaction.Status = "completed"
	if err := db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to complete transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction completed successfully",
		Data:    transaction,
	})
}

// CancelTransaction - Cancel transaction
func CancelTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var transaction models.Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, TransactionResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && transaction.UserID != userID {
		c.JSON(http.StatusForbidden, TransactionResponse{
			Success: false,
			Message: "You don't have permission to cancel this transaction",
		})
		return
	}

	if transaction.Status == "completed" || transaction.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, TransactionResponse{
			Success: false,
			Message: "Transaction cannot be cancelled",
		})
		return
	}

	// Reverse stock for OUT transactions
	if transaction.Type == "OUT" && transaction.Status != "cancelled" {
		var items []models.TransactionItem
		db.Where("transaction_id = ?", transaction.ID).Find(&items)

		for _, item := range items {
			var stock models.Stock
			if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err == nil {
				stock.Quantity += item.Quantity
				stock.LastUpdated = time.Now()
				db.Save(&stock)
			}
		}
	}

	transaction.Status = "cancelled"
	if err := db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, TransactionResponse{
			Success: false,
			Message: "Failed to cancel transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Message: "Transaction cancelled successfully",
		Data:    transaction,
	})
}

// GetTransactionStats - Get transaction statistics
func GetTransactionStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var stats struct {
		Total     int64   `json:"total"`
		Draft     int64   `json:"draft"`
		Pending   int64   `json:"pending"`
		Approved  int64   `json:"approved"`
		Completed int64   `json:"completed"`
		Cancelled int64   `json:"cancelled"`
		TotalIN   int64   `json:"total_in"`
		TotalOUT  int64   `json:"total_out"`
		Revenue   float64 `json:"revenue"`
	}

	query := db.Model(&models.Transaction{})

	if role != "super_admin" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// Count by status
	query.Count(&stats.Total)
	query.Where("status = ?", "draft").Count(&stats.Draft)
	query.Where("status = ?", "pending").Count(&stats.Pending)
	query.Where("status = ?", "approved").Count(&stats.Approved)
	query.Where("status = ?", "completed").Count(&stats.Completed)
	query.Where("status = ?", "cancelled").Count(&stats.Cancelled)

	// Count by type
	query.Where("type = ? AND status = ?", "IN", "completed").Count(&stats.TotalIN)
	query.Where("type = ? AND status = ?", "OUT", "completed").Count(&stats.TotalOUT)

	// Total revenue (only completed IN transactions)
	db.Model(&models.Transaction{}).
		Where("type = ? AND status = ?", "IN", "completed").
		Select("COALESCE(SUM(total_value), 0)").
		Scan(&stats.Revenue)

	c.JSON(http.StatusOK, TransactionResponse{
		Success: true,
		Data:    stats,
	})
}