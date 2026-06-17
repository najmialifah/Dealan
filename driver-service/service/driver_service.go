package service

import (
	"context"
	"errors"

	"github.com/najmialifah/Dealan/driver-service/domain"
	"github.com/najmialifah/Dealan/driver-service/repository"
)

type driverService struct {
	repo repository.DriverRepository
}

// NewDriverService membuat instance baru bisnis service driver
func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{
		repo: repo,
	}
}

// UpdateLocation memperbarui koordinat GPS driver ke database
func (s *driverService) UpdateLocation(ctx context.Context, id string, req domain.UpdateLocationRequest) error {
	return s.repo.UpdateDriverLocation(ctx, id, req.Lat, req.Long)
}

// UpdateStatus mengubah status online/offline dan layanan aktif pengemudi
func (s *driverService) UpdateStatus(ctx context.Context, id string, req domain.UpdateStatusRequest) error {
	return s.repo.UpdateDriverStatus(ctx, id, req.IsOnline, req.LayananAktif)
}

// GetProfile mengambil profil lengkap pengemudi
func (s *driverService) GetProfile(ctx context.Context, id string) (*domain.Driver, error) {
	return s.repo.GetDriverByID(ctx, id)
}

// CreateDriver mendaftarkan driver baru yang dipicu oleh event pendaftaran dari Kafka
func (s *driverService) CreateDriver(ctx context.Context, id, name, email, phone string) error {
	existing, _ := s.repo.GetDriverByID(ctx, id)
	if existing != nil {
		return nil // idempotent
	}

	driver := &domain.Driver{
		ID:      id,
		Nama:    name,
		Email:   email,
		NomorHP: phone,
		Status:  "active",
	}
	return s.repo.CreateDriver(ctx, driver)
}

// AddVehicle menambahkan kendaraan baru ke akun driver
func (s *driverService) AddVehicle(ctx context.Context, driverID string, vehicle domain.Vehicle) error {
	driver, err := s.repo.GetDriverByID(ctx, driverID)
	if err != nil {
		return err
	}
	if driver == nil {
		return errors.New("driver tidak ditemukan")
	}

	vehicle.DriverID = driverID
	driver.Vehicles = append(driver.Vehicles, vehicle)
	return s.repo.UpdateDriver(ctx, driver)
}
