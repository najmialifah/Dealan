package service

import (
	"errors"
	"location-service/domain"
)

type LocationService struct {
	Repo domain.LocationRepository
}

func (s *LocationService) UpdateDriverLocation(loc *domain.LocationUpdate) error {
	if loc.EntityID == "" {
		return errors.New("entity ID tidak boleh kosong")
	}
	if loc.Latitude < -90 || loc.Latitude > 90 {
		return errors.New("latitude tidak valid")
	}
	if loc.Longitude < -180 || loc.Longitude > 180 {
		return errors.New("longitude tidak valid")
	}
	return s.Repo.SaveLocation(loc)
}