package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware — fixed version
// Menghindari kombinasi terlarang Access-Control-Allow-Origin: * dengan credentials
func CORSMiddleware() gin.HandlerFunc {
	// Daftar origin yang diizinkan (development + production)
	allowedOrigins := map[string]bool{
		"http://localhost:4200": true, // Angular dev
		"http://127.0.0.1:4200": true,
		"http://localhost:4300": true, // port alternatif
		"http://localhost:5173": true, // Vite
		"http://localhost:3000": true, // React/Next
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// ✅ Set origin SPESIFIK (bukan wildcard) saat credentials aktif
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Vary", "Origin") // penting untuk caching
		} else if origin == "" {
			// Request tanpa Origin header (curl, Postman) → izinkan wildcard
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		// Kalau origin tidak diizinkan → jangan set header → browser block

		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, OPTIONS, GET, PUT, PATCH, DELETE")

		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept, Accept-Encoding, "+
				"X-CSRF-Token, Authorization, Origin, Cache-Control, X-Requested-With")

		c.Writer.Header().Set("Access-Control-Expose-Headers",
			"Content-Length, Content-Disposition")

		c.Writer.Header().Set("Access-Control-Max-Age", "43200") // 12 jam

		// Handle preflight (OPTIONS)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}