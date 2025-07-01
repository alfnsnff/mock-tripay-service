package models

import "time"

// Request Models
type CreateTransactionRequest struct {
	Method        string      `json:"method" binding:"required"`
	MerchantRef   string      `json:"merchant_ref" binding:"required"`
	Amount        int         `json:"amount" binding:"required,min=1"`
	CustomerName  string      `json:"customer_name" binding:"required"`
	CustomerEmail string      `json:"customer_email" binding:"required,email"`
	CustomerPhone string      `json:"customer_phone" binding:"required"`
	OrderItems    []OrderItem `json:"order_items"`
	ReturnURL     string      `json:"return_url"`
	ExpiredTime   int64       `json:"expired_time"`
	Signature     string      `json:"signature"`
}

type OrderItem struct {
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	ProductURL string `json:"product_url,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
}

// Response Models
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type TransactionData struct {
	Reference            string      `json:"reference"`
	MerchantRef          string      `json:"merchant_ref"`
	PaymentSelectionType string      `json:"payment_selection_type"`
	PaymentMethod        string      `json:"payment_method"`
	PaymentName          string      `json:"payment_name"`
	CustomerName         string      `json:"customer_name"`
	CustomerEmail        string      `json:"customer_email"`
	CustomerPhone        string      `json:"customer_phone"`
	CallbackURL          string      `json:"callback_url"`
	ReturnURL            string      `json:"return_url"`
	Amount               int         `json:"amount"`
	Fee                  int         `json:"fee"`
	TotalFee             int         `json:"total_fee"`
	AmountReceived       int         `json:"amount_received"`
	PayCode              string      `json:"pay_code"`
	PayURL               string      `json:"pay_url"`
	CheckoutURL          string      `json:"checkout_url"`
	Status               string      `json:"status"`
	ExpiredTime          int64       `json:"expired_time"`
	OrderItems           []OrderItem `json:"order_items"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type TransactionDetailData struct {
	Reference      string     `json:"reference"`
	MerchantRef    string     `json:"merchant_ref"`
	PaymentMethod  string     `json:"payment_method"`
	PaymentName    string     `json:"payment_name"`
	CustomerName   string     `json:"customer_name"`
	CustomerEmail  string     `json:"customer_email"`
	CustomerPhone  string     `json:"customer_phone"`
	Amount         int        `json:"amount"`
	Fee            int        `json:"fee"`
	TotalFee       int        `json:"total_fee"`
	AmountReceived int        `json:"amount_received"`
	Status         string     `json:"status"`
	PaidAt         *time.Time `json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type PaymentChannel struct {
	Group       string                 `json:"group"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	FeeMerchant map[string]interface{} `json:"fee_merchant"`
	FeeCustomer map[string]interface{} `json:"fee_customer"`
	TotalFee    map[string]interface{} `json:"total_fee"`
	MinimumFee  int                    `json:"minimum_fee"`
	MaximumFee  int                    `json:"maximum_fee"`
	IconURL     string                 `json:"icon_url"`
	Active      bool                   `json:"active"`
}

// Statistics for monitoring
type MockStats struct {
	TotalRequests     int     `json:"total_requests"`
	SuccessRequests   int     `json:"success_requests"`
	FailureRequests   int     `json:"failure_requests"`
	SuccessRate       float64 `json:"success_rate"`
	TotalTransactions int     `json:"total_transactions"`
	PaidTransactions  int     `json:"paid_transactions"`
	AvgResponseTime   int64   `json:"avg_response_time_ms"`
}
