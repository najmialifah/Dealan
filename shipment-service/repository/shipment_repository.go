package repository

import (
	"context"

	"github.com/najmialifah/Dealan/shipment-service/domain"
	"gorm.io/gorm"
)

// ShipmentRepository mendefinisikan kontrak GORM database untuk CRUD Shipment.
type ShipmentRepository interface {
	Create(ctx context.Context, shipment *domain.Shipment) error
	GetByID(ctx context.Context, id string) (*domain.Shipment, error)
	GetByTrackingCode(ctx context.Context, code string) (*domain.Shipment, error)
	Update(ctx context.Context, shipment *domain.Shipment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Shipment, error)
}

type shipmentRepositoryImpl struct {
	db *gorm.DB
}

// NewShipmentRepository membuat instance baru dari ShipmentRepository.
func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepositoryImpl{db: db}
}

func (r *shipmentRepositoryImpl) Create(ctx context.Context, shipment *domain.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

func (r *shipmentRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepositoryImpl) GetByTrackingCode(ctx context.Context, code string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	err := r.db.WithContext(ctx).Where("tracking_code = ?", code).First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepositoryImpl) Update(ctx context.Context, shipment *domain.Shipment) error {
	return r.db.WithContext(ctx).Save(shipment).Error
}

func (r *shipmentRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Shipment{}).Error
}

func (r *shipmentRepositoryImpl) List(ctx context.Context) ([]domain.Shipment, error) {
	var shipments []domain.Shipment
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&shipments).Error
	if err != nil {
		return nil, err
	}
	return shipments, nil
}
