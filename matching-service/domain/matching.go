package domain

import (
	"context"
)

// MatchRequest menangkap data JSON dari Controller
type MatchRequest struct {
	OrderID   uint    `json:"order_id" binding:"required"` // Menggunakan uint sesuai kebutuhan service kamu
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
	Radius    int     `json:"radius"` // Default nanti diatur di service (5000)
}

// MatchedDriver adalah balikan dari Database / Repository
type MatchedDriver struct {
	DriverID  string  `json:"driver_id"` // Pakai string karena UUID di Supabase
	Distance  float64 `json:"distance"`  // Jarak dalam meter
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ========================================================
// KONTRAK INTERFACE (Satu-satunya sumber kebenaran)
// ========================================================

// MatchingRepository adalah kontrak untuk file repository
type MatchingRepository interface {
	FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters int) (*MatchedDriver, error)
}

// MatchingService adalah kontrak untuk file service
type MatchingService interface {
	MatchOrder(ctx context.Context, req *MatchRequest) (*MatchedDriver, error)
}