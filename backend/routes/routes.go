package routes

import (
	"github.com/gin-gonic/gin"

	"atk-backend/handlers"
	"atk-backend/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// ============================================================
	// HANDLERS
	// ============================================================
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	roleHandler := handlers.NewRoleHandler()
	masterHandler := handlers.NewMasterHandler()
	stockHandler := handlers.NewStockHandler()
	transactionHandler := handlers.NewTransactionHandler()
	documentHandler := handlers.NewDocumentHandler()
	mockHandler := handlers.NewMockHandler()

	// ============================================================
	// PUBLIC ROUTES
	// ============================================================
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// ============================================================
	// PROTECTED ROUTES
	// ============================================================
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// ---- Auth ----
		api.GET("/auth/validate", authHandler.ValidateToken)

		// ============================================================
		// PROFILE
		// ============================================================
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)
		api.POST("/profile/change-password", userHandler.ChangePassword)

		// ============================================================
		// USERS
		// ============================================================
		api.GET("/users", userHandler.GetAll)
		api.POST("/users", userHandler.Create)
		api.POST("/users/import", userHandler.ImportExcel)
		api.GET("/users/export", userHandler.ExportExcel)
		api.GET("/users/username/:username", userHandler.GetByUsername)
		api.GET("/users/:id", userHandler.GetByID)
		api.PUT("/users/:id", userHandler.Update)
		api.DELETE("/users/:id", userHandler.Delete)
		api.POST("/users/:id/roles", userHandler.AssignRole)
		api.DELETE("/users/:id/roles", userHandler.RemoveRole)
		api.GET("/users/:id/roles", userHandler.GetUserRoles)

		// ============================================================
		// ROLES
		// ============================================================
		api.GET("/roles", roleHandler.GetAll)
		api.POST("/roles", roleHandler.Create)
		api.POST("/roles/import", roleHandler.ImportExcel)
		api.GET("/roles/export", roleHandler.ExportExcel)
		api.GET("/roles/:id", roleHandler.GetByID)
		api.PUT("/roles/:id", roleHandler.Update)
		api.DELETE("/roles/:id", roleHandler.Delete)

		// ============================================================
		// MASTER DATA
		// ============================================================
		api.GET("/master", masterHandler.GetAll)
		api.POST("/master", masterHandler.Create)
		api.POST("/master/import", masterHandler.ImportExcel)
		api.GET("/master/export", masterHandler.ExportExcel)
		api.GET("/master/template", masterHandler.DownloadTemplate)
		api.GET("/master/:id", masterHandler.GetByID)
		api.PUT("/master/:id", masterHandler.Update)
		api.DELETE("/master/:id", masterHandler.Delete)

		// Alias lama
		api.GET("/masters", masterHandler.GetAll)
		api.POST("/masters", masterHandler.Create)

		// ============================================================
		// STOCK
		// ============================================================
		api.GET("/stock", stockHandler.GetAll)
		api.POST("/stock", stockHandler.Create)
		api.POST("/stock/import", stockHandler.ImportExcel)
		api.GET("/stock/export", stockHandler.ExportExcel)
		api.GET("/stock/:id", stockHandler.GetByID)
		api.PUT("/stock/:id", stockHandler.Update)
		api.DELETE("/stock/:id", stockHandler.Delete)

		// Alias lama /api/stocks
		api.GET("/stocks", stockHandler.GetAll)
		api.GET("/stocks/master/:masterId", stockHandler.GetByMasterID)
		api.PUT("/stocks/:id", stockHandler.Update)
		api.GET("/stocks/low-stock", stockHandler.GetLowStock)
		api.GET("/stocks/summary", stockHandler.GetStockSummary)
		api.POST("/stocks/adjust/:masterId", stockHandler.AdjustStock)

		// ============================================================
		// TRANSACTIONS — ✅ SEMUA PAKAI HANDLER ASLI
		// ============================================================
		api.GET("/transactions", transactionHandler.GetAll)
		api.POST("/transactions", transactionHandler.Create)
		api.POST("/transactions/import", transactionHandler.ImportExcel)
		api.GET("/transactions/export", transactionHandler.ExportExcel)
		api.GET("/transactions/:id", transactionHandler.GetByID)
		api.PUT("/transactions/:id", transactionHandler.Update)
		api.DELETE("/transactions/:id", transactionHandler.Delete)
		api.POST("/transactions/:id/submit", transactionHandler.Submit)
		api.POST("/transactions/:id/approve", transactionHandler.Approve)
		api.POST("/transactions/:id/reject", transactionHandler.Reject)
		api.POST("/transactions/:id/complete", transactionHandler.Complete)
		api.POST("/transactions/:id/cancel", transactionHandler.Cancel)

		// ============================================================
		// APPROVALS (mock)
		// ============================================================
		api.GET("/approvals", mockHandler.Paged)
		api.POST("/approvals/:id/approve", mockHandler.Ok)
		api.POST("/approvals/:id/reject", mockHandler.Ok)
		api.DELETE("/approvals/:id", mockHandler.Ok)

		// ============================================================
		// DASHBOARD (mock)
		// ============================================================
		api.GET("/dashboard/summary", mockHandler.DashboardSummary)
		api.GET("/dashboard/trend", mockHandler.DashboardTrend)
		api.GET("/dashboard/stock-by-category", mockHandler.DashboardStockByCategory)

		// ============================================================
		// REPORTS (mock)
		// ============================================================
		api.GET("/reports/stock", mockHandler.ReportList)
		api.GET("/reports/transaction", mockHandler.ReportList)
		api.GET("/reports/approval", mockHandler.ReportList)
		api.GET("/reports/stock/export", mockHandler.Export)
		api.GET("/reports/transaction/export", mockHandler.Export)
		api.GET("/reports/approval/export", mockHandler.Export)

		// ============================================================
		// DOCUMENTS — ✅ VALIDATE PAKAI HANDLER ASLI
		// ============================================================
		api.GET("/documents", documentHandler.GetAll)
		api.POST("/documents", documentHandler.Create)
		api.GET("/documents/code/:code", documentHandler.GetByCode)
		api.POST("/documents/scan", documentHandler.ScanQRCode)
		api.GET("/documents/validate", documentHandler.ValidateDocument)   // ✅ GANTI DARI mockHandler
		api.GET("/documents/:id", documentHandler.GetByID)
		api.POST("/documents/:id/sign", documentHandler.SignDocument)
		api.GET("/documents/:id/qr", documentHandler.GenerateQRCode)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ATK Enterprise API is running",
		})
	})

	return router
}