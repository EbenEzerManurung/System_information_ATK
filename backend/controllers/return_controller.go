package controllers

import (
	"net/http"
	"strconv"
	"time"

	"atk-backend/database"
	"atk-backend/models"

	"github.com/gin-gonic/gin"
)

type ReturnResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// CreateReturn - Create new return transaction
func CreateReturn(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	var input struct {
		TransactionID uint   `json:"transaction_id" binding:"required"`
		Items         []struct {
			MasterID uint `json:"master_id" binding:"required"`
			Quantity int  `json:"quantity" binding:"required,gt=0"`
			Reason   string `json:"reason"`
		} `json:"items" binding:"required,min=1"`
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Check if transaction exists
	var transaction models.Transaction
	if err := db.Preload("Items").First(&transaction, input.TransactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	// Check if transaction is completed
	if transaction.Status != "completed" {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Only completed transactions can be returned",
		})
		return
	}

	// Check if user has permission (staff can return own transactions, admin can return all)
	if role != "super_admin" && role != "admin" && transaction.UserID != userID {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "You don't have permission to return this transaction",
		})
		return
	}

	// Generate return code
	code := "RET-" + time.Now().Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%10000, 10)

	// Create return transaction
	returnTrans := models.Transaction{
		TransactionCode: code,
		Type:            "RETURN",
		Status:          "draft",
		UserID:          userID,
		Notes:           input.Notes,
		TransactionDate: time.Now(),
	}

	var totalItems int
	var totalValue float64
	var returnItems []models.TransactionItem

	// Process each return item
	for _, item := range input.Items {
		// Check if item exists in original transaction
		var foundItem models.TransactionItem
		found := false
		for _, txItem := range transaction.Items {
			if txItem.MasterID == item.MasterID {
				foundItem = txItem
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusBadRequest, ReturnResponse{
				Success: false,
				Message: "Item not found in original transaction",
			})
			return
		}

		// Check if quantity is valid (cannot return more than original)
		if item.Quantity > foundItem.Quantity {
			c.JSON(http.StatusBadRequest, ReturnResponse{
				Success: false,
				Message: "Cannot return more than original quantity",
			})
			return
		}

		// Get master details
		var master models.Master
		if err := db.First(&master, item.MasterID).Error; err != nil {
			c.JSON(http.StatusNotFound, ReturnResponse{
				Success: false,
				Message: "Master item not found",
			})
			return
		}

		// Calculate subtotal
		subTotal := float64(item.Quantity) * master.Price
		totalItems += item.Quantity
		totalValue += subTotal

		// Create return item
		returnItem := models.TransactionItem{
			MasterID: item.MasterID,
			Quantity: item.Quantity,
			Price:    master.Price,
			SubTotal: subTotal,
			Notes:    item.Reason,
		}
		returnItems = append(returnItems, returnItem)
	}

	returnTrans.TotalItems = totalItems
	returnTrans.TotalValue = totalValue
	returnTrans.Items = returnItems

	// Save return transaction
	if err := db.Create(&returnTrans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to create return: " + err.Error(),
		})
		return
	}

	// Update stock (add back quantity)
	for _, item := range returnItems {
		var stock models.Stock
		if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err == nil {
			stock.Quantity += item.Quantity
			stock.LastUpdated = time.Now()
			db.Save(&stock)
		}
	}

	c.JSON(http.StatusCreated, ReturnResponse{
		Success: true,
		Message: "Return created successfully",
		Data:    returnTrans,
	})
}

// GetReturns - Get all return transactions
func GetReturns(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	offset := (page - 1) * limit

	query := db.Model(&models.Transaction{}).
		Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master").
		Where("type = ?", "RETURN")

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// If not admin, only show own returns
	if role != "super_admin" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Get returns
	var returns []models.Transaction
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to get returns: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Returns retrieved successfully",
		Data:    returns,
		Total:   total,
	})
}

// GetReturnByID - Get return by ID
func GetReturnByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var returnTrans models.Transaction
	if err := db.Preload("User").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Master").
		Where("id = ? AND type = ?", id, "RETURN").
		First(&returnTrans).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Return not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && returnTrans.UserID != userID {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "You don't have permission to view this return",
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Return retrieved successfully",
		Data:    returnTrans,
	})
}

// ApproveReturn - Approve return transaction
func ApproveReturn(c *gin.Context) {
	id := c.Param("id")
	approverID := c.GetUint("user_id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "Only admin or manager can approve returns",
		})
		return
	}

	db := database.GetDB()

	var returnTrans models.Transaction
	if err := db.Where("id = ? AND type = ?", id, "RETURN").First(&returnTrans).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Return not found",
		})
		return
	}

	if returnTrans.Status != "pending" {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Return is not pending",
		})
		return
	}

	// Update status
	returnTrans.Status = "approved"
	returnTrans.ApprovedBy = &approverID
	now := time.Now()
	returnTrans.ApprovedAt = &now

	if err := db.Save(&returnTrans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to approve return: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Return approved successfully",
		Data:    returnTrans,
	})
}

// SubmitReturn - Submit return for approval
func SubmitReturn(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var returnTrans models.Transaction
	if err := db.Where("id = ? AND type = ?", id, "RETURN").First(&returnTrans).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Return not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && returnTrans.UserID != userID {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "You don't have permission to submit this return",
		})
		return
	}

	if returnTrans.Status != "draft" {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Only draft returns can be submitted",
		})
		return
	}

	returnTrans.Status = "pending"
	if err := db.Save(&returnTrans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to submit return: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Return submitted successfully",
		Data:    returnTrans,
	})
}

// CancelReturn - Cancel return transaction
func CancelReturn(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var returnTrans models.Transaction
	if err := db.Where("id = ? AND type = ?", id, "RETURN").First(&returnTrans).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Return not found",
		})
		return
	}

	// Check permission
	if role != "super_admin" && role != "admin" && returnTrans.UserID != userID {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "You don't have permission to cancel this return",
		})
		return
	}

	if returnTrans.Status == "completed" || returnTrans.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Return cannot be cancelled",
		})
		return
	}

	// Reverse stock if approved
	if returnTrans.Status == "approved" {
		for _, item := range returnTrans.Items {
			var stock models.Stock
			if err := db.Where("master_id = ?", item.MasterID).First(&stock).Error; err == nil {
				stock.Quantity -= item.Quantity
				stock.LastUpdated = time.Now()
				db.Save(&stock)
			}
		}
	}

	returnTrans.Status = "cancelled"
	if err := db.Save(&returnTrans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to cancel return: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Return cancelled successfully",
		Data:    returnTrans,
	})
}

// CompleteReturn - Complete return transaction
func CompleteReturn(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "super_admin" && role != "admin" {
		c.JSON(http.StatusForbidden, ReturnResponse{
			Success: false,
			Message: "Only admin can complete returns",
		})
		return
	}

	db := database.GetDB()

	var returnTrans models.Transaction
	if err := db.Where("id = ? AND type = ?", id, "RETURN").First(&returnTrans).Error; err != nil {
		c.JSON(http.StatusNotFound, ReturnResponse{
			Success: false,
			Message: "Return not found",
		})
		return
	}

	if returnTrans.Status != "approved" {
		c.JSON(http.StatusBadRequest, ReturnResponse{
			Success: false,
			Message: "Return must be approved first",
		})
		return
	}

	returnTrans.Status = "completed"
	if err := db.Save(&returnTrans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to complete return: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Return completed successfully",
		Data:    returnTrans,
	})
}

// GetReturnByTransaction - Get returns by original transaction ID
func GetReturnByTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("role")

	db := database.GetDB()

	var returns []models.Transaction
	query := db.Preload("User").
		Preload("Items").
		Preload("Items.Master").
		Where("type = ? AND id IN (SELECT return_transaction_id FROM transaction_returns WHERE original_transaction_id = ?)", "RETURN", transactionID)

	if role != "super_admin" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ReturnResponse{
			Success: false,
			Message: "Failed to get returns: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReturnResponse{
		Success: true,
		Message: "Returns retrieved successfully",
		Data:    returns,
		Total:   int64(len(returns)),
	})
}