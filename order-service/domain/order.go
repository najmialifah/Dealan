package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Order merepresentasikan tabel orders di database
type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      string         `gorm:"not null" json:"user_id"`
	DriverID    *string        `json:"driver_id"` // Bisa null jika belum ada driver
	Status      string         `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	DetailPaket map[string]any `gorm:"type:jsonb" json:"detail_paket"` // Menggunakan tipe data JSONB di PostgreSQL
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateOrderRequest adalah payload request untuk membuat order
type CreateOrderRequest struct {
	UserID      string         `json:"user_id" binding:"required"`
	DetailPaket map[string]any `json:"detail_paket" binding:"required"`
}

// OrderRepository merupakan antarmuka (interface) untuk interaksi database order
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, id uint) (*Order, error)
}

// OrderService mendefinisikan logika bisnis terkait order
type OrderService interface {
	CreateOrder(ctx context.Context, order *Order) error
}