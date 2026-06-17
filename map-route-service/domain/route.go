package domain

import (
	"context"
	"time"
)

// MapRoute mendefinisikan skema penyimpanan rute jalan dengan polyline bertipe TEXT di PostgreSQL.
type MapRoute struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Origin      string    `gorm:"type:varchar(255);not null;index:idx_origin_dest" json:"origin"`
	Destination string    `gorm:"type:varchar(255);not null;index:idx_origin_dest" json:"destination"`
	Polyline    string    `gorm:"type:text;not null" json:"polyline"`           // Rute jalan yang diserialisasi (Polyline)
	Distance    float64   `gorm:"type:decimal(10,2);not null" json:"distance"` // Jarak rute dalam km
	Duration    int       `gorm:"type:integer;not null" json:"duration"`       // Durasi rute dalam detik
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RouteRequest adalah payload untuk meminta data rute baru.
type RouteRequest struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// RouteResponse mengembalikan rute jalan lengkap beserta visualisasi garis polylinenya.
type RouteResponse struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Polyline    string  `json:"polyline"`
	Distance    float64 `json:"distance"` // km
	Duration    int     `json:"duration"` // detik
}

// MapRepository mendefinisikan interface operasi database untuk rute jalan.
type MapRepository interface {
	GetRoute(ctx context.Context, origin, destination string) (*MapRoute, error)
	SaveRoute(ctx context.Context, route *MapRoute) error
}

// MapService menyediakan kontrak fungsional untuk map-route-service.
type MapService interface {
	GetOrCreateRoute(ctx context.Context, req RouteRequest) (*RouteResponse, error)
}
