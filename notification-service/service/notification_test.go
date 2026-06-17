package service_test

import (
	"context"
	"testing"

	"github.com/najmialifah/Dealan/notification-service/domain"
	"github.com/najmialifah/Dealan/notification-service/mocks"
	"github.com/najmialifah/Dealan/notification-service/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationService_SendNotification_Success(t *testing.T) {
	// 1. Inisialisasi controller mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Buat repository palsu (Mock)
	mockRepo := mocks.NewMockNotificationRepository(ctrl)

	// 3. Masukkan mock ke service
	svc := service.NewNotificationService(mockRepo)

	req := domain.NotificationRequest{
		TargetID:   "user-123",
		Title:      "Info Promo",
		Body:       "Diskon 50% untuk perjalanan selanjutnya!",
		ActionLink: "/promo",
	}

	// Ekspektasikan bahwa SaveLog dipanggil 1 kali dengan argumen apa pun untuk response dan mengembalikan nil (tanpa error)
	mockRepo.EXPECT().
		SaveLog(gomock.Any(), req, gomock.Any()).
		Return(nil).
		Times(1)

	// 4. Jalankan fungsi
	resp, err := svc.SendNotification(context.Background(), req)

	// 5. Verifikasi hasil menggunakan testify
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Sent", resp.DeliveryStatus)
	assert.NotEmpty(t, resp.MessageID)
}