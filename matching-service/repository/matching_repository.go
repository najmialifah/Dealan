package repository

import (
	"context"
	"matching-service/models"

	"gorm.io/gorm"
)

type MatchingRepository interface {
	FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters float64) (*models.MatchedDriver, error)
}

type matchingRepository struct {
	db *gorm.DB
}

func NewMatchingRepository(db *gorm.DB) MatchingRepository {
	return &matchingRepository{db}
}

// FindNearestDriver menggunakan Raw SQL dengan fungsi spasial PostGIS ST_DWithin
// untuk mencari 1 driver terdekat dari koordinat order.
func (r *matchingRepository) FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters float64) (*models.MatchedDriver, error) {
	var driver models.MatchedDriver

	// Asumsi tabel driver_locations di-share atau direplikasi ke database matching.
	// Jika berbeda DB, maka arsitektur harus memfasilitasi integrasi (misal via Kafka sink).
	// Di sini kita asumsikan matching service memiliki akses read ke tabel yang sama atau sinkron.
	query := `
		SELECT 
			driver_id, 
			ST_Y(location::geometry) as latitude, 
			ST_X(location::geometry) as longitude,
			ST_Distance(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) as distance
		FROM driver_locations
		WHERE ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		ORDER BY distance ASC
		LIMIT 1;
	`

	// Ingat bahwa ST_MakePoint menerima parameter (Longitude, Latitude) -> (X, Y)
	err := r.db.WithContext(ctx).Raw(query, lon, lat, lon, lat, radiusMeters).Scan(&driver).Error
	if err != nil {
		return nil, err
	}

	// Jika tidak ada data yang ditemukan, struct akan kosong (driver_id = 0)
	// Kita bisa melakukan handle error gorm.ErrRecordNotFound jika menggunakan First,
	// tetapi karena Raw().Scan(), jika row kosong maka tidak return error secara otomatis, 
	// kita cek secara manual dari driver_id.
	if driver.DriverID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &driver, nil
}
