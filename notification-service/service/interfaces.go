package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/notification-service/domain"
)

// AuthService interface mendefinisikan fungsi apa saja yang tersedia di service ini
type NotificationService interface {
	SendNotification(ctx context.Context, req domain.NotificationRequest) (*domain.NotificationResponse, error)
}