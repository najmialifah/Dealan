package service_test

import (
	"context"
	"testing"

	"github.com/shakilaaulia/Dealan/notification-service/domain"
	"github.com/shakilaaulia/Dealan/notification-service/mocks" 
	"github.com/shakilaaulia/Dealan/notification-service/service"
	"go.uber.org/mock/gomock"
)

func TestNotificationService_Send_ExpectedFail(t *testing.T) {
	// 1. Inisialisasi controller mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Buat repository palsu (Mock)
	mockRepo := mocks.NewMockNotificationRepository(ctrl)

	// 3. Masukkan mock ke service
	svc := service.NewNotificationService(mockRepo)

	req := domain.NotificationRequest{
		TargetID: "U1",
		   Title: "Promo"}
	
	// 4. Jalankan fungsi. 
	// Karena di notification_service.go kita tulis "not implemented", maka resp akan nil.
	resp, err := svc.SendNotification(context.Background(), req)

	// 5. Cek apakah benar-benar gagal
	if err == nil {
		t.Errorf("Harusnya error 'not implemented', tapi malah sukses")
	}
	if resp != nil {
		t.Errorf("Harusnya response nil karena gagal, tapi malah ada isinya")
	}
}