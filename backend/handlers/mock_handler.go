package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockHandler — untuk endpoint yang belum diimplementasikan
type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

// ============================================================
// GENERIC RESPONSES
// ============================================================

// Ok — response sukses generik untuk POST/PUT/DELETE
func (h *MockHandler) Ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK (mock)",
		"data":    nil,
	})
}

// Paged — response paginasi kosong untuk GET list
func (h *MockHandler) Paged(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data": gin.H{
			"data":       []gin.H{},
			"total":      0,
			"page":       1,
			"size":       10,
			"totalPages": 0,
		},
	})
}

// Export — response file Excel kosong
func (h *MockHandler) Export(c *gin.Context) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=export.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", []byte{})
}

// ============================================================
// DASHBOARD
// ============================================================

func (h *MockHandler) DashboardSummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data": gin.H{
			"totalUsers":           3,
			"totalItems":           0,
			"monthTransactions":    0,
			"pendingApprovals":     0,
			"recentTransactions":   []gin.H{},
			"lowStock":             []gin.H{},
			"pendingApprovalsList": []gin.H{},
		},
	})
}

func (h *MockHandler) DashboardTrend(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    []gin.H{},
	})
}

func (h *MockHandler) DashboardStockByCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    []gin.H{},
	})
}

// ============================================================
// REPORT
// ============================================================

func (h *MockHandler) ReportList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data": gin.H{
			"data":       []gin.H{},
			"total":      0,
			"page":       1,
			"size":       10,
			"totalPages": 0,
			"summary":    []gin.H{},
		},
	})
}

// ============================================================
// DOCUMENT
// ============================================================

func (h *MockHandler) DocumentValidate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data": gin.H{
			"valid":    false,
			"document": nil,
		},
	})
}