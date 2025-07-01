package main

import (
	"log"
	"os"

	"mock-tripay/internal/handlers"
	"mock-tripay/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "3001")
	mode := getEnv("GIN_MODE", "release")

	// Set Gin mode
	gin.SetMode(mode)

	// Initialize services
	paymentService := services.NewPaymentService()

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup router
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(requestCounterMiddleware())

	// Routes
	setupRoutes(router, paymentHandler)

	// Start server
	log.Printf("🎭 Mock Tripay Server starting on port %s", port)
	log.Printf("📋 Mode: %s", mode)
	printAvailableEndpoints()

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, handler *handlers.PaymentHandler) {
	// Health check
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Tripay-compatible endpoints
		api.GET("/merchant/payment-channel", handler.GetPaymentChannels)
		api.POST("/transaction/create", handler.CreateTransaction)
		api.GET("/transaction/detail", handler.GetTransactionDetail) // Changed from /:reference to query param

		// Mock server utilities
		api.GET("/stats", handler.GetStats)
		api.POST("/reset", handler.ResetData)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func requestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple request logging
		log.Printf("📥 %s %s - %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func printAvailableEndpoints() {
	log.Println("📋 Available endpoints:")
	log.Println("   GET  /health")
	log.Println("   GET  /api/merchant/payment-channel")
	log.Println("   POST /api/transaction/create")
	log.Println("   GET  /api/transaction/detail?reference=<reference>") // Updated documentation
	log.Println("   GET  /api/stats")
	log.Println("   POST /api/reset")
	log.Println("")
	log.Println("🔧 To use with eTicket API:")
	log.Println("   Set TRIPAY_BASE_URL=http://localhost:3001/api")
}
