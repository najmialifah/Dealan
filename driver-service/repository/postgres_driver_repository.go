package repository

import (
	"context"
	"time"

	"github.com/najmialifah/Dealan/driver-service/domain"
	"gorm.io/gorm"
)

type postgresDriverRepository struct {
	db *gorm.DB
}

// NewPostgresDriverRepository membuat instance baru DriverRepository berbasis GORM PostgreSQL
func NewPostgresDriverRepository(db *gorm.DB) DriverRepository {
	return &postgresDriverRepository{
		db: db,
	}
}

// GetDriverByID mengambil data driver lengkap beserta relasi vehicle, status, dan rating
func (r *postgresDriverRepository) GetDriverByID(ctx context.Context, id string) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.db.WithContext(ctx).
		Preload("Vehicles").
		Preload("DriverStatus").
		Preload("DriverRating").
		Where("id = ?", id).
		First(&driver).Error
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

// UpdateDriver menyimpan pembaruan model driver ke database
func (r *postgresDriverRepository) UpdateDriver(ctx context.Context, driver *domain.Driver) error {
	return r.db.WithContext(ctx).Save(driver).Error
}

// CreateDriver mendaftarkan driver baru beserta status dan rating default-nya dalam satu transaksi ACID
func (r *postgresDriverRepository) CreateDriver(ctx context.Context, driver *domain.Driver) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Simpan driver utama
		if err := tx.Create(driver).Error; err != nil {
			return err
		}

		// Buat status online default
		status := &domain.DriverStatus{
			DriverID:     driver.ID,
			IsOnline:     false,
			LayananAktif: "",
			LastSeen:     time.Now(),
			Lat:          0.0,
			Long:         0.0,
		}
		if err := tx.Create(status).Error; err != nil {
			return err
		}

		// Buat rating default
		rating := &domain.DriverRating{
			DriverID:    driver.ID,
			TotalRating: 5.00,
			TotalReview: 0,
			UpdatedAt:   time.Now(),
		}
		if err := tx.Create(rating).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateDriverLocation memperbarui koordinat GPS lokasi pengemudi secara efisien
func (r *postgresDriverRepository) UpdateDriverLocation(ctx context.Context, driverID string, lat, long float64) error {
	return r.db.WithContext(ctx).Model(&domain.DriverStatus{}).
		Where("driver_id = ?", driverID).
		Updates(map[string]interface{}{
			"lat":       lat,
			"long":      long,
			"last_seen": time.Now(),
		}).Error
}

// UpdateDriverStatus memperbarui status online/offline dan jenis layanan aktif
func (r *postgresDriverRepository) UpdateDriverStatus(ctx context.Context, driverID string, isOnline bool, layananAktif string) error {
	return r.db.WithContext(ctx).Model(&domain.DriverStatus{}).
		Where("driver_id = ?", driverID).
		Updates(map[string]interface{}{
			"is_online":     isOnline,
			"layanan_aktif": layananAktif,
			"last_seen":     time.Now(),
		}).Error
}
