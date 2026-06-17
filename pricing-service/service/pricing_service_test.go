package service

import (
	"context"
	"errors"
	"testing"

	"github.com/najmialifah/Dealan/pricing-service/models"
	"github.com/stretchr/testify/assert"
)

// mockPricingRepository adalah implementasi stub dari repository.PricingRepository untuk kebutuhan unit test.
type mockPricingRepository struct {
	rule *models.PricingRule
	err  error
}

func (m *mockPricingRepository) GetRuleByServiceType(ctx context.Context, serviceType string) (*models.PricingRule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rule, nil
}

func (m *mockPricingRepository) SaveRule(ctx context.Context, rule *models.PricingRule) error {
	return nil
}

func (m *mockPricingRepository) SaveNegotiation(ctx context.Context, neg *models.PricingNegotiation) error {
	return nil
}

func (m *mockPricingRepository) GetNegotiationByOrderID(ctx context.Context, orderID string) (*models.PricingNegotiation, error) {
	return nil, nil
}

// TestCalculateEstimate menguji perhitungan tarif estimasi dinamis.
func TestCalculateEstimate(t *testing.T) {
	mockRule := &models.PricingRule{
		ServiceType: "ride",
		BasePrice:   8000.0,
		Active:      true,
		Config: models.JSONB{
			"per_km_rate":           2500.0,
			"rush_hour_multiplier":  1.5,
			"night_multiplier":      1.2,
			"min_price":             10000.0,
			"negotiation_tolerance": 0.20,
		},
	}

	repo := &mockPricingRepository{rule: mockRule}
	svc := NewPricingService(repo)

	t.Run("✅ Tarif normal dengan jarak dekat (< 2 km)", func(t *testing.T) {
		req := models.PricingEstimateRequest{
			ServiceType: "ride",
			Distance:    1.5,
		}
		res, err := svc.CalculateEstimate(context.Background(), req)
		assert.NoError(t, err)
		// Harus bernilai min_price karena hasil kalkulasi (8000) kurang dari min_price (10000)
		assert.Equal(t, 10000.0, res.EstimatedPrice)
	})

	t.Run("✅ Tarif normal dengan jarak jauh (> 2 km)", func(t *testing.T) {
		req := models.PricingEstimateRequest{
			ServiceType: "ride",
			Distance:    4.0, // 8000 + (2.0 * 2500) = 13000
		}
		res, err := svc.CalculateEstimate(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 13000.0, res.EstimatedPrice)
		assert.Equal(t, 10400.0, res.MinPrice) // 13000 - 20%
		assert.Equal(t, 15600.0, res.MaxPrice) // 13000 + 20%
	})

	t.Run("✅ Tarif jam sibuk (Rush Hour)", func(t *testing.T) {
		req := models.PricingEstimateRequest{
			ServiceType: "ride",
			Distance:    4.0, // 13000
			IsRushHour:  true, // 13000 * 1.5 = 19500
		}
		res, err := svc.CalculateEstimate(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 19500.0, res.EstimatedPrice)
	})

	t.Run("❌ Gagal jika rule tidak ditemukan", func(t *testing.T) {
		errRepo := &mockPricingRepository{err: errors.New("database error")}
		errSvc := NewPricingService(errRepo)

		req := models.PricingEstimateRequest{
			ServiceType: "ride",
			Distance:    4.0,
		}
		_, err := errSvc.CalculateEstimate(context.Background(), req)
		assert.Error(t, err)
	})
}

// TestNegotiatePrice menguji validasi status penawaran harga.
func TestNegotiatePrice(t *testing.T) {
	repo := &mockPricingRepository{}
	svc := NewPricingService(repo)

	t.Run("✅ Penawaran disetujui (dalam batas toleransi)", func(t *testing.T) {
		req := models.NegotiationRequest{
			OrderID:        "ORD-001",
			OriginalPrice:  15000.0,
			RequestedPrice: 13000.0, // Masih di atas batas minimal 15000 * 0.8 = 12000
		}
		res, err := svc.NegotiatePrice(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "approved", res.Status)
	})

	t.Run("❌ Penawaran ditolak (di luar batas toleransi)", func(t *testing.T) {
		req := models.NegotiationRequest{
			OrderID:        "ORD-002",
			OriginalPrice:  15000.0,
			RequestedPrice: 10000.0, // Di bawah batas minimal 12000
		}
		res, err := svc.NegotiatePrice(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "rejected", res.Status)
	})
}
