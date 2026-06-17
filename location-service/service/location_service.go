package service

import (
	"context"
	"location-service/domain"
	"location-service/repository"
	"log"
)

// LocationService mendefinisikan logika bisnis untuk location-service
type LocationService interface {
	UpdateLocation(ctx context.Context, driverID uint, lat, lon float64) error
	GetNearbyDrivers(ctx context.Context, lat, lon float64, radius float64) ([]domain.NearbyDriver, error)
}

type locationService struct {
	repo repository.LocationRepository
}

// NewLocationService menginisialisasi location service
func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{
		repo: repo,
	}
}

// UpdateLocation memperbarui lokasi secara asynchronous menggunakan Goroutine
func (s *locationService) UpdateLocation(ctx context.Context, driverID uint, lat, lon float64) error {
	// Menggunakan Goroutine agar update koordinat tidak ngeblok HTTP response.
	// Kita bisa memanfaatkan background context agar goroutine tidak dibatalkan saat http request selesai.
	bgCtx := context.Background()
	
	go func(bgContext context.Context, dID uint, latitude, longitude float64) {
		err := s.repo.UpsertLocation(bgContext, dID, latitude, longitude)
		if err != nil {
			// Di sistem produksi yang sesungguhnya kita dapat push metrics/log alert
			log.Printf("[Error] Gagal update lokasi asinkron driver %d: %v", dID, err)
		}
	}(bgCtx, driverID, lat, lon)

	// Kembalikan langsung (fire and forget) untuk merespon API lebih cepat
	return nil
}

// GetNearbyDrivers mengambil daftar driver yang ada di dalam radius tertentu
func (s *locationService) GetNearbyDrivers(ctx context.Context, lat, lon float64, radius float64) ([]domain.NearbyDriver, error) {
	// Panggil repository untuk query ST_DWithin
	return s.repo.FindNearbyDrivers(ctx, lat, lon, radius)
}