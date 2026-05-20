package service

import (
	"matching-service/domain"
	"matching-service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestFindDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMatchingRepository(ctrl)
	matchService := MatchingService{Repo: mockRepo}

	t.Run("Gagal Order ID Kosong", func(t *testing.T) {
		req := &domain.MatchingRequest{OrderID: "", ServiceType: "ride"}
		_, err := matchService.FindDriver(req)
		if err == nil {
			t.Errorf("Harusnya error karena Order ID kosong")
		}
	})

	t.Run("Sukses Cari Driver", func(t *testing.T) {
		req := &domain.MatchingRequest{OrderID: "ORD001", ServiceType: "ride"}
		resp := &domain.MatchingResponse{DriverID: "DRV123", EstimasiWaktu: 5}

		// Ekspektasi mock balikin data driver
		mockRepo.EXPECT().FindNearestDriver(req).Return(resp, nil).Times(1)

		hasil, err := matchService.FindDriver(req)
		if err != nil || hasil.DriverID != "DRV123" {
			t.Errorf("Gagal mendapatkan driver yang sesuai")
		}
	})
}