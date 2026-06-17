package service

import (
	"context"

	"github.com/najmialifah/Dealan/driver-service/domain"
)

// DriverService mendefinisikan interface bisnis untuk driver-service
type DriverService interface {
	UpdateLocation(ctx context.Context, id string, req domain.UpdateLocationRequest) error
	UpdateStatus(ctx context.Context, id string, req domain.UpdateStatusRequest) error
	GetProfile(ctx context.Context, id string) (*domain.Driver, error)
	CreateDriver(ctx context.Context, id, name, email, phone string) error
	AddVehicle(ctx context.Context, driverID string, vehicle domain.Vehicle) error
}
