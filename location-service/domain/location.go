package domain

import (
	"context"
	"time"
)

type LocationUpdate struct {
	DriverID  uint    `json:"driver_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Bearing   float64 `json:"bearing,omitempty"`
}

type DriverLocation struct {
	DriverID  uint      `json:"driver_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NearbyDriver struct {
	DriverID  uint    `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance_meters"`
}

type LocationRepository interface {
	UpsertLocation(ctx context.Context, driverID uint, lat, lon float64) error
	FindNearbyDrivers(ctx context.Context, lat, lon float64, radiusMeters float64) ([]NearbyDriver, error)
}

type LocationService interface {
	UpdateLocation(ctx context.Context, driverID uint, lat, lon float64) error
	GetNearbyDrivers(ctx context.Context, lat, lon float64, radius float64) ([]NearbyDriver, error)
}