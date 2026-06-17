package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/najmialifah/Dealan/notification-service/domain"
	"github.com/najmialifah/Dealan/notification-service/repository"
)

type notificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService membuat instance baru untuk layanan notifikasi
func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

// SendNotification mensimulasikan pengiriman notifikasi ke FCM/APNS dan menyimpan log ke DB
func (s *notificationService) SendNotification(ctx context.Context, req domain.NotificationRequest) (*domain.NotificationResponse, error) {
	// Generate message ID unik sebagai simulasi pengiriman sukses
	messageID := fmt.Sprintf("msg-%d-%d", time.Now().UnixNano(), rand.Intn(1000))

	res := domain.NotificationResponse{
		DeliveryStatus: "Sent",
		MessageID:      messageID,
	}

	// Simpan log ke database PostgreSQL menggunakan repository
	if s.repo != nil {
		err := s.repo.SaveLog(ctx, req, res)
		if err != nil {
			fmt.Printf("[Notification Service] Gagal menyimpan log notifikasi ke DB: %v\n", err)
		}
	}

	return &res, nil
}