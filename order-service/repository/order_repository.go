package repository

import (
	"context"
	"order-service/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository menginisialisasi repository order
func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db}
}

// CreateOrder menyimpan order baru ke dalam database
func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetOrderByID mengambil order berdasarkan ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
