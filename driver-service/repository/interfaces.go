package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mock_driver_repository.go -package=mocks

type DriverRepository interface {
	GetDriverByID(ctx context.Context, id string) (*domain.Driver, error)
	UpdateDriver(ctx context.Context, driver *domain.Driver) error
}
