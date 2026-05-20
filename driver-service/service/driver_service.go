package service

import (
	"context"
	"errors"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
	"github.com/shakilaaulia/Dealan/driver-service/repository"
)

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{
		repo: repo,
	}
}

func (s *driverService) UpdateLocation(ctx context.Context, id string, req domain.UpdateLocationRequest) error {
	driver, err := s.repo.GetDriverByID(ctx, id)
	if err != nil {
		return err
	}
	if driver == nil {
		return errors.New("driver not found")
	}

	driver.Lat = req.Lat
	driver.Long = req.Long
	return s.repo.UpdateDriver(ctx, driver)
}

func (s *driverService) UpdateStatus(ctx context.Context, id string, req domain.UpdateStatusRequest) error {
	driver, err := s.repo.GetDriverByID(ctx, id)
	if err != nil {
		return err
	}
	if driver == nil {
		return errors.New("driver not found")
	}

	driver.Status = req.Status
	return s.repo.UpdateDriver(ctx, driver)
}

func (s *driverService) GetProfile(ctx context.Context, id string) (*domain.Driver, error) {
	return s.repo.GetDriverByID(ctx, id)
}
