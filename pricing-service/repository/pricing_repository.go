package repository

import (
	"context"

	"github.com/najmialifah/Dealan/pricing-service/models"
	"gorm.io/gorm"
)

// PricingRepository menyediakan interface operasi database untuk pricing-service.
type PricingRepository interface {
	GetRuleByServiceType(ctx context.Context, serviceType string) (*models.PricingRule, error)
	SaveRule(ctx context.Context, rule *models.PricingRule) error
	SaveNegotiation(ctx context.Context, neg *models.PricingNegotiation) error
	GetNegotiationByOrderID(ctx context.Context, orderID string) (*models.PricingNegotiation, error)
}

type pricingRepositoryImpl struct {
	db *gorm.DB
}

// NewPricingRepository membuat instance baru dari PricingRepository.
func NewPricingRepository(db *gorm.DB) PricingRepository {
	return &pricingRepositoryImpl{db: db}
}

// GetRuleByServiceType mengambil aturan harga aktif berdasarkan jenis layanan.
func (r *pricingRepositoryImpl) GetRuleByServiceType(ctx context.Context, serviceType string) (*models.PricingRule, error) {
	var rule models.PricingRule
	err := r.db.WithContext(ctx).Where("service_type = ? AND active = ?", serviceType, true).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// SaveRule menyimpan atau memperbarui aturan harga di database.
func (r *pricingRepositoryImpl) SaveRule(ctx context.Context, rule *models.PricingRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// SaveNegotiation menyimpan data negosiasi harga ke database.
func (r *pricingRepositoryImpl) SaveNegotiation(ctx context.Context, neg *models.PricingNegotiation) error {
	return r.db.WithContext(ctx).Save(neg).Error
}

// GetNegotiationByOrderID mengambil riwayat negosiasi berdasarkan ID pesanan.
func (r *pricingRepositoryImpl) GetNegotiationByOrderID(ctx context.Context, orderID string) (*models.PricingNegotiation, error) {
	var neg models.PricingNegotiation
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at desc").First(&neg).Error
	if err != nil {
		return nil, err
	}
	return &neg, nil
}
