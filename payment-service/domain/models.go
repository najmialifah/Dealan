package domain

type PaymentRequest struct {
	OrderID          string  `json:"order_id"`
	Nominal          float64 `json:"nominal"`
	MetodePembayaran string  `json:"metode_pembayaran"` // Cash, QRIS, Bank
	PromoID          string  `json:"promo_id"`
}

type PaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentStatus string `json:"payment_status"`
	InvoiceURL    string `json:"invoice_url"`
}
