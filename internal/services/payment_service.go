package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"mock-tripay/internal/models"

	"github.com/google/uuid"
)

type PaymentService struct {
	transactions map[string]*models.TransactionDetailData
	mutex        sync.RWMutex
	stats        *models.MockStats

	// Configuration
	successRate     float32
	avgResponseTime time.Duration
	autoPayDelay    time.Duration
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		transactions:    make(map[string]*models.TransactionDetailData),
		mutex:           sync.RWMutex{},
		stats:           &models.MockStats{},
		successRate:     0.999, // 98% success rate
		avgResponseTime: 100 * time.Millisecond,
		autoPayDelay:    30 * time.Second, // Auto-pay after 30 seconds
	}
}

func (s *PaymentService) CreateTransaction(req *models.CreateTransactionRequest) (*models.APIResponse, error) {
	start := time.Now()
	defer func() {
		s.updateStats(time.Since(start), true)
	}()

	// Simulate processing delay
	s.simulateDelay()

	// Simulate random failure
	if rand.Float32() > s.successRate {
		s.updateStats(time.Since(start), false)
		return &models.APIResponse{
			Success: false,
			Message: "Payment gateway temporarily unavailable",
		}, nil
	}

	// Generate unique reference
	reference := s.generateReference()

	// Calculate fees
	fee := s.calculateFee(req.Amount, req.Method)
	amountReceived := req.Amount - fee

	// Create transaction record
	transactionDetail := &models.TransactionDetailData{
		Reference:      reference,
		MerchantRef:    req.MerchantRef,
		PaymentMethod:  req.Method,
		PaymentName:    s.getPaymentName(req.Method),
		CustomerName:   req.CustomerName,
		CustomerEmail:  req.CustomerEmail,
		CustomerPhone:  req.CustomerPhone,
		Amount:         req.Amount,
		Fee:            fee,
		TotalFee:       fee,
		AmountReceived: amountReceived,
		Status:         "UNPAID",
		PaidAt:         nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store transaction
	s.mutex.Lock()
	s.transactions[reference] = transactionDetail
	s.stats.TotalTransactions++
	s.mutex.Unlock()

	// Start auto-payment simulation in background
	go s.simulateAutoPayment(reference)

	// Build response data
	responseData := &models.TransactionData{
		Reference:            reference,
		MerchantRef:          req.MerchantRef,
		PaymentSelectionType: "static",
		PaymentMethod:        req.Method,
		PaymentName:          s.getPaymentName(req.Method),
		CustomerName:         req.CustomerName,
		CustomerEmail:        req.CustomerEmail,
		CustomerPhone:        req.CustomerPhone,
		CallbackURL:          "",
		ReturnURL:            req.ReturnURL,
		Amount:               req.Amount,
		Fee:                  fee,
		TotalFee:             fee,
		AmountReceived:       amountReceived,
		PayCode:              s.generatePayCode(req.Method),
		PayURL:               fmt.Sprintf("mock://tripay.com/pay/%s", reference),
		CheckoutURL:          fmt.Sprintf("mock://tripay.com/checkout/%s", reference),
		Status:               "UNPAID",
		ExpiredTime:          req.ExpiredTime,
		OrderItems:           req.OrderItems,
		CreatedAt:            transactionDetail.CreatedAt,
		UpdatedAt:            transactionDetail.UpdatedAt,
	}

	return &models.APIResponse{
		Success: true,
		Message: "Transaction created",
		Data:    responseData,
	}, nil
}

func (s *PaymentService) GetTransactionDetail(reference string) (*models.APIResponse, error) {
	start := time.Now()
	defer func() {
		s.updateStats(time.Since(start), true)
	}()

	s.simulateDelay()

	s.mutex.RLock()
	transaction, exists := s.transactions[reference]
	s.mutex.RUnlock()

	if !exists {
		return &models.APIResponse{
			Success: false,
			Message: "Transaction not found",
		}, nil
	}

	return &models.APIResponse{
		Success: true,
		Message: "Get transaction success",
		Data:    transaction,
	}, nil
}

func (s *PaymentService) GetPaymentChannels() *models.APIResponse {
	s.simulateDelay()

	channels := []models.PaymentChannel{
		{
			Group: "Virtual Account",
			Code:  "BRIVA",
			Name:  "BRI Virtual Account",
			Type:  "virtual_account",
			FeeMerchant: map[string]interface{}{
				"flat":    4000,
				"percent": 0,
			},
			FeeCustomer: map[string]interface{}{
				"flat":    0,
				"percent": 0,
			},
			TotalFee: map[string]interface{}{
				"flat":    4000,
				"percent": 0,
			},
			MinimumFee: 0,
			MaximumFee: 0,
			IconURL:    "https://tripay.co.id/images/payment_icon/bri.png",
			Active:     true,
		},
		{
			Group: "E-Wallet",
			Code:  "QRIS",
			Name:  "QRIS",
			Type:  "qr_code",
			FeeMerchant: map[string]interface{}{
				"flat":    0,
				"percent": 0.7,
			},
			FeeCustomer: map[string]interface{}{
				"flat":    0,
				"percent": 0,
			},
			TotalFee: map[string]interface{}{
				"flat":    0,
				"percent": 0.7,
			},
			MinimumFee: 0,
			MaximumFee: 0,
			IconURL:    "https://tripay.co.id/images/payment_icon/qris.png",
			Active:     true,
		},
		{
			Group: "Convenience Store",
			Code:  "ALFAMART",
			Name:  "Alfamart",
			Type:  "convenience_store",
			FeeMerchant: map[string]interface{}{
				"flat":    2500,
				"percent": 0,
			},
			FeeCustomer: map[string]interface{}{
				"flat":    0,
				"percent": 0,
			},
			TotalFee: map[string]interface{}{
				"flat":    2500,
				"percent": 0,
			},
			MinimumFee: 0,
			MaximumFee: 0,
			IconURL:    "https://tripay.co.id/images/payment_icon/alfamart.png",
			Active:     true,
		},
	}

	return &models.APIResponse{
		Success: true,
		Message: "Get payment channels success",
		Data:    channels,
	}
}

func (s *PaymentService) GetStats() *models.APIResponse {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Calculate success rate
	if s.stats.TotalRequests > 0 {
		s.stats.SuccessRate = float64(s.stats.SuccessRequests) / float64(s.stats.TotalRequests)
	}

	return &models.APIResponse{
		Success: true,
		Message: "Mock server statistics",
		Data:    s.stats,
	}
}

func (s *PaymentService) ResetData() *models.APIResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.transactions = make(map[string]*models.TransactionDetailData)
	s.stats = &models.MockStats{}

	return &models.APIResponse{
		Success: true,
		Message: "Mock data reset successfully",
	}
}

// Helper methods
func (s *PaymentService) simulateDelay() {
	// Add randomness to response time (50ms - 150ms)
	delay := s.avgResponseTime + time.Duration(rand.Intn(20))*time.Millisecond
	time.Sleep(delay)
}

func (s *PaymentService) generateReference() string {
	return fmt.Sprintf("T%d%s", time.Now().Unix(), uuid.New().String()[:8])
}

func (s *PaymentService) calculateFee(amount int, method string) int {
	switch method {
	case "QRIS":
		return int(float64(amount) * 0.007) // 0.7%
	case "BRIVA", "BCAVA", "MANDIRIVA":
		return 4000 // Flat fee
	case "ALFAMART", "INDOMARET":
		return 2500 // Flat fee
	default:
		return int(float64(amount) * 0.01) // 1%
	}
}

func (s *PaymentService) getPaymentName(method string) string {
	paymentNames := map[string]string{
		"QRIS":      "QRIS",
		"BRIVA":     "BRI Virtual Account",
		"BCAVA":     "BCA Virtual Account",
		"MANDIRIVA": "Mandiri Virtual Account",
		"ALFAMART":  "Alfamart",
		"INDOMARET": "Indomaret",
	}

	if name, exists := paymentNames[method]; exists {
		return name
	}
	return "Unknown Payment Method"
}

func (s *PaymentService) generatePayCode(method string) string {
	switch method {
	case "QRIS":
		return "" // QRIS doesn't have pay code
	case "BRIVA", "BCAVA", "MANDIRIVA":
		return fmt.Sprintf("%d", rand.Intn(9000000000)+1000000000) // 10 digits
	case "ALFAMART", "INDOMARET":
		return fmt.Sprintf("%d", rand.Intn(900000000)+100000000) // 9 digits
	default:
		return fmt.Sprintf("%d", rand.Intn(9000)+1000) // 4 digits
	}
}

func (s *PaymentService) simulateAutoPayment(reference string) {
	// Wait for auto-pay delay
	time.Sleep(s.autoPayDelay)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	transaction, exists := s.transactions[reference]
	if !exists || transaction.Status != "UNPAID" {
		return
	}

	// Complete payment
	now := time.Now()
	transaction.Status = "PAID"
	transaction.PaidAt = &now
	transaction.UpdatedAt = now

	s.stats.PaidTransactions++
}

func (s *PaymentService) updateStats(duration time.Duration, success bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats.TotalRequests++
	if success {
		s.stats.SuccessRequests++
	} else {
		s.stats.FailureRequests++
	}

	// Update average response time
	s.stats.AvgResponseTime = (s.stats.AvgResponseTime + duration.Milliseconds()) / 2
}
