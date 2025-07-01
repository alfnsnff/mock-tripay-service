package handlers

import (
	"log"
	"net/http"

	"mock-tripay/internal/models"
	"mock-tripay/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// POST /api/transaction/create
func (h *PaymentHandler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.service.CreateTransaction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/transaction/detail?reference=<reference>
func (h *PaymentHandler) GetTransactionDetail(c *gin.Context) {
	reference := c.Query("reference")

	// Debug logging
	log.Printf("🔍 GetTransactionDetail called")
	log.Printf("   - URL: %s", c.Request.URL.String())
	log.Printf("   - Reference from query: '%s'", reference)
	log.Printf("   - All query params: %v", c.Request.URL.Query())

	if reference == "" {
		log.Printf("❌ Empty reference parameter")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Payment reference is required",
		})
		return
	}

	log.Printf("🔍 Calling service.GetTransactionDetail with: %s", reference)

	response, err := h.service.GetTransactionDetail(reference)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	log.Printf("📤 Service response: Success=%t, Message=%s", response.Success, response.Message)

	if !response.Success {
		log.Printf("❌ Service returned success=false: %s", response.Message)
		c.JSON(http.StatusNotFound, response)
		return
	}

	log.Printf("✅ Returning successful response")
	c.JSON(http.StatusOK, response)
}

// GET /api/merchant/payment-channel
func (h *PaymentHandler) GetPaymentChannels(c *gin.Context) {
	response := h.service.GetPaymentChannels()
	c.JSON(http.StatusOK, response)
}

// GET /api/stats
func (h *PaymentHandler) GetStats(c *gin.Context) {
	response := h.service.GetStats()
	c.JSON(http.StatusOK, response)
}

// POST /api/reset
func (h *PaymentHandler) ResetData(c *gin.Context) {
	response := h.service.ResetData()
	c.JSON(http.StatusOK, response)
}

// GET /health
func (h *PaymentHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "Mock Tripay server is running",
		"timestamp": c.Request.Header.Get("X-Request-Time"),
		"service":   "mock-tripay",
		"version":   "1.0.0",
	})
}
