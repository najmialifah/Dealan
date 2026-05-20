package repository

import (
	"context"
	"github.com/shakilaaulia/Dealan/notification-service/domain"
)

// NotificationRepository adalah interface (kontrak) untuk bagian database.
// Ini digunakan oleh gomock untuk membuat 'database palsu' saat testing.
type NotificationRepository interface {
	SaveLog(ctx context.Context, req domain.NotificationRequest, res domain.NotificationResponse) error
}