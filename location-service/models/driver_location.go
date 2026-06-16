package models

import (
	"time"
)

// DriverLocation merepresentasikan lokasi terbaru driver
type DriverLocation struct {
	DriverID  uint      `gorm:"primaryKey" json:"driver_id"`
	Latitude  float64   `gorm:"-" json:"latitude"`  // Tidak masuk ke skema kolom tabel langsung, digabung di Geom
	Longitude float64   `gorm:"-" json:"longitude"` // Tidak masuk ke skema kolom tabel langsung, digabung di Geom
	Location  string    `gorm:"type:geometry(Point,4326)" json:"-"` // PostGIS geometry Point (Longitude, Latitude)
	UpdatedAt time.Time `json:"updated_at"`
}

// NearbyDriver response struct
type NearbyDriver struct {
	DriverID  uint    `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance_meters"`
}
