package domain

type Order struct {
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
}

type OrderRepository interface {
	Save(order *Order) error
}