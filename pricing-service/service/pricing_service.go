package service

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/shakilaaulia/Dealan/pricing-service/models"
	"github.com/shakilaaulia/Dealan/pricing-service/repository"
)

// PricingService mendefinisikan kontrak fungsi layanan kalkulasi harga dan negosiasi.
type PricingService interface {
	CalculateEstimate(ctx context.Context, req models.PricingEstimateRequest) (models.PricingEstimateResponse, error)
	NegotiatePrice(ctx context.Context, req models.NegotiationRequest) (models.NegotiationResponse, error)
}

type pricingServiceImpl struct {
	repo repository.PricingRepository
	// Kafka writer atau logger bisa ditaruh di sini jika dibutuhkan.
}

// NewPricingService membuat instance baru untuk layanan kalkulasi harga.
func NewPricingService(repo repository.PricingRepository) PricingService {
	return &pricingServiceImpl{repo: repo}
}

// CalculateEstimate mengkalkulasi tarif estimasi berdasarkan aturan tarif dinamis yang disimpan di JSONB.
func (s *pricingServiceImpl) CalculateEstimate(ctx context.Context, req models.PricingEstimateRequest) (models.PricingEstimateResponse, error) {
	rule, err := s.repo.GetRuleByServiceType(ctx, req.ServiceType)
	if err != nil {
		return models.PricingEstimateResponse{}, fmt.Errorf("gagal mendapatkan aturan tarif: %v", err)
	}

	// Ekstrasi data konfigurasi tarif dari JSONB dengan nilai default yang aman
	perKmRate := s.getFloat64(rule.Config, "per_km_rate", 2500.0)
	rushHourMultiplier := s.getFloat64(rule.Config, "rush_hour_multiplier", 1.5)
	nightMultiplier := s.getFloat64(rule.Config, "night_multiplier", 1.2)
	minPrice := s.getFloat64(rule.Config, "min_price", 10000.0)
	weightRatePerKg := s.getFloat64(rule.Config, "weight_rate_per_kg", 1000.0)
	negotiationTolerance := s.getFloat64(rule.Config, "negotiation_tolerance", 0.20) // Toleransi negosiasi ±20%

	// Kalkulasi dasar: Tarif awal mencakup 2 km pertama
	basePrice := rule.BasePrice
	estimatedPrice := basePrice

	if req.Distance > 2.0 {
		estimatedPrice += (req.Distance - 2.0) * perKmRate
	}

	// Tambahan tarif khusus layanan pengiriman barang (GoSend) berdasarkan berat
	if req.ServiceType == "send" && req.Weight > 0 {
		estimatedPrice += req.Weight * weightRatePerKg
	}

	// Terapkan faktor pengali jam sibuk (Rush Hour)
	if req.IsRushHour {
		estimatedPrice *= rushHourMultiplier
	}

	// Terapkan faktor pengali perjalanan malam hari
	if req.IsNight {
		estimatedPrice *= nightMultiplier
	}

	// Pastikan harga tidak berada di bawah harga minimum
	if estimatedPrice < minPrice {
		estimatedPrice = minPrice
	}

	// Bulatkan harga ke kelipatan 500 terdekat agar lebih rapi bagi pelanggan
	estimatedPrice = math.Round(estimatedPrice/500) * 500

	// Hitung batas harga minimal dan maksimal untuk fitur negosiasi
	minNegotiationPrice := estimatedPrice * (1.0 - negotiationTolerance)
	maxNegotiationPrice := estimatedPrice * (1.0 + negotiationTolerance)

	// Bulatkan batas negosiasi
	minNegotiationPrice = math.Round(minNegotiationPrice/500) * 500
	maxNegotiationPrice = math.Round(maxNegotiationPrice/500) * 500

	return models.PricingEstimateResponse{
		ServiceType:    req.ServiceType,
		Distance:       req.Distance,
		EstimatedPrice: estimatedPrice,
		MinPrice:       minNegotiationPrice,
		MaxPrice:       maxNegotiationPrice,
	}, nil
}

// NegotiatePrice memvalidasi penawaran harga dari pengguna berdasarkan batas toleransi.
func (s *pricingServiceImpl) NegotiatePrice(ctx context.Context, req models.NegotiationRequest) (models.NegotiationResponse, error) {
	// Di sistem nyata, kita mencari tarif estimasi asli di database atau state pesanan.
	// Di sini kita memvalidasi penawaran harga: toleransi standar adalah ±20% dari harga asli.
	tolerance := 0.20
	minAllowed := req.OriginalPrice * (1.0 - tolerance)
	maxAllowed := req.OriginalPrice * (1.0 + tolerance)

	status := "approved"
	if req.RequestedPrice < minAllowed || req.RequestedPrice > maxAllowed {
		status = "rejected"
	}

	// Simpan riwayat negosiasi ke database
	negotiation := &models.PricingNegotiation{
		OrderID:        req.OrderID,
		OriginalPrice:  req.OriginalPrice,
		RequestedPrice: req.RequestedPrice,
		Status:         status,
	}

	if err := s.repo.SaveNegotiation(ctx, negotiation); err != nil {
		log.Printf("[PricingService] Gagal menyimpan riwayat negosiasi: %v", err)
	}

	return models.NegotiationResponse{
		OrderID:        req.OrderID,
		OriginalPrice:  req.OriginalPrice,
		RequestedPrice: req.RequestedPrice,
		Status:         status,
	}, nil
}

// Helper untuk membaca float64 dari JSONB secara aman
func (s *pricingServiceImpl) getFloat64(config models.JSONB, key string, defaultValue float64) float64 {
	if val, ok := config[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return defaultValue
}
