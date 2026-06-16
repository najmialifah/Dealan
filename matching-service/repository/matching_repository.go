package repository // WAJIB package repository

import (
	"context"
	"matching-service/domain"

	"gorm.io/gorm"
)

type matchingRepository struct {
	db *gorm.DB
}

func NewMatchingRepository(db *gorm.DB) domain.MatchingRepository {
	return &matchingRepository{db}
}

func (r *matchingRepository) FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters int) (*domain.MatchedDriver, error) {
	var driver domain.MatchedDriver

	query := `
		SELECT 
			driver_id, 
			ST_Y(lokasi::geometry) as latitude, 
			ST_X(lokasi::geometry) as longitude,
			ST_Distance(lokasi::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) as distance
		FROM driver_status
		WHERE is_online = true
		AND ST_DWithin(lokasi::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		ORDER BY distance ASC
		LIMIT 1;
	`

	err := r.db.WithContext(ctx).Raw(query, lon, lat, lon, lat, radiusMeters).Scan(&driver).Error
	if err != nil {
		return nil, err
	}

	if driver.DriverID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &driver, nil
}