package repository

import (
	"context"

	"github.com/najmialifah/Dealan/driver-service/domain"
)

// DriverRepository mendefinisikan interface database untuk driver-service menggunakan GORM
type DriverRepository interface {
	GetDriverByID(ctx context.Context, id string) (*domain.Driver, error)
	UpdateDriver(ctx context.Context, driver *domain.Driver) error
	CreateDriver(ctx context.Context, driver *domain.Driver) error
	UpdateDriverLocation(ctx context.Context, driverID string, lat, long float64) error
	UpdateDriverStatus(ctx context.Context, driverID string, isOnline bool, layananAktif string) error
}
