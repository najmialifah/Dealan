package service

import (
	"context"
	"map-route-service/domain"
)

type MapService interface {
	GetRoute(ctx context.Context, req domain.RouteRequest) (*domain.RouteResponse, error)
}