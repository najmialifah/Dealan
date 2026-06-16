package repository

import (
	"context"
	"order-service/models"

	"gorm.io/gorm"
)

// OrderRepository merupakan antarmuka (interface) untuk interaksi database order
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository menginisialisasi repository order
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

// CreateOrder menyimpan order baru ke dalam database
func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetOrderByID mengambil order berdasarkan ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
