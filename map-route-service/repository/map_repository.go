package repository

import (
	"context"

	"map-route-service/domain"
	"gorm.io/gorm"
)

type mapRepositoryImpl struct {
	db *gorm.DB
}

// NewMapRepository membuat instance baru untuk MapRepository menggunakan database GORM.
func NewMapRepository(db *gorm.DB) domain.MapRepository {
	return &mapRepositoryImpl{db: db}
}

// GetRoute mengambil rute yang sudah tersimpan sebelumnya berdasarkan origin dan destination.
func (r *mapRepositoryImpl) GetRoute(ctx context.Context, origin, destination string) (*domain.MapRoute, error) {
	var route domain.MapRoute
	err := r.db.WithContext(ctx).Where("origin = ? AND destination = ?", origin, destination).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

// SaveRoute menyimpan data rute jalan baru ke database PostgreSQL.
func (r *mapRepositoryImpl) SaveRoute(ctx context.Context, route *domain.MapRoute) error {
	return r.db.WithContext(ctx).Save(route).Error
}
