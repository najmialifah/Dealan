package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
)

type DriverService interface {
	UpdateLocation(ctx context.Context, id string, req domain.UpdateLocationRequest) error
	UpdateStatus(ctx context.Context, id string, req domain.UpdateStatusRequest) error
	GetProfile(ctx context.Context, id string) (*domain.Driver, error)
}
