package service

import (
	"context"
	"map-route-service/domain"
	"map-route-service/repository"
)

type mapService struct {
	repo repository.MapRepository
}

func NewMapService(r repository.MapRepository) MapService {
	return &mapService{r}
}

func (s *mapService) GetRoute(ctx context.Context, req domain.RouteRequest) (*domain.RouteResponse, error) {
	return s.repo.GetRoute(ctx, req.Origin, req.Destination)
}