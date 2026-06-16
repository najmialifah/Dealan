package repository

import (
	"context"
	"location-service/models"

	"gorm.io/gorm"
)

// LocationRepository merupakan antarmuka untuk operasi ke database PostgreSQL/PostGIS
type LocationRepository interface {
	UpsertLocation(ctx context.Context, driverID uint, lat, lon float64) error
	FindNearbyDrivers(ctx context.Context, lat, lon float64, radiusMeters float64) ([]models.NearbyDriver, error)
}

type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepository menginisialisasi repository location
func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db}
}

// UpsertLocation memperbarui lokasi driver, menggunakan raw SQL PostGIS karena kompleksitas query upsert geom
func (r *locationRepository) UpsertLocation(ctx context.Context, driverID uint, lat, lon float64) error {
	// Memasukkan/memperbarui titik geometri. PostGIS menggunakan koordinat X Y (Longitude, Latitude)
	query := `
		INSERT INTO driver_locations (driver_id, location, updated_at) 
		VALUES (?, ST_SetSRID(ST_MakePoint(?, ?), 4326), NOW())
		ON CONFLICT (driver_id) 
		DO UPDATE SET location = EXCLUDED.location, updated_at = NOW();
	`
	return r.db.WithContext(ctx).Exec(query, driverID, lon, lat).Error
}

// FindNearbyDrivers menggunakan fungsi PostGIS ST_DWithin untuk mencari koordinat di dalam radius
func (r *locationRepository) FindNearbyDrivers(ctx context.Context, lat, lon float64, radiusMeters float64) ([]models.NearbyDriver, error) {
	var nearbyDrivers []models.NearbyDriver

	// Query spasial ST_DWithin: menggunakan fungsi geografi agar perhitungan jarak akurat dalam satuan meter
	// Cast ke geography diperlukan jika kolom tipe dasar geometry tanpa parameter geographic
	query := `
		SELECT 
			driver_id, 
			ST_Y(location::geometry) as latitude, 
			ST_X(location::geometry) as longitude,
			ST_Distance(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) as distance
		FROM driver_locations
		WHERE ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		ORDER BY distance ASC
		LIMIT 50;
	`

	// Parameter untuk query adalah (lon, lat)
	err := r.db.WithContext(ctx).Raw(query, lon, lat, lon, lat, radiusMeters).Scan(&nearbyDrivers).Error
	if err != nil {
		return nil, err
	}

	return nearbyDrivers, nil
}
