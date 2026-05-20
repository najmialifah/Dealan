package repository

import (
	"context"
	"map-route-service/domain"
)

type MapRepository interface {
	GetRoute(ctx context.Context, origin string, destination string) (*domain.RouteResponse, error)
}