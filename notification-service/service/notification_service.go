package service

import (
	"context"
	"errors"

	"github.com/shakilaaulia/Dealan/notification-service/domain"
	"github.com/shakilaaulia/Dealan/notification-service/repository"
)

// notificationService adalah implementasi nyata dari interface NotificationService
type notificationService struct {
	repo repository.NotificationRepository // Bergantung pada repository (Dependency Injection)
}

// fungsi untuk membuat instance service baru
func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

// logika pengiriman notifikasi
func (s *notificationService) SendNotification(ctx context.Context, req domain.NotificationRequest) (*domain.NotificationResponse, error) {

	return nil, errors.New("not implemented")
}