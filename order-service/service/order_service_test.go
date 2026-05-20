package service

import (
	"order-service/domain"
	"order-service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCreateOrder(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := OrderService{Repo: mockRepo}

	t.Run("Gagal karena tipe layanan kosong", func(t *testing.T) {
		orderKosong := &domain.Order{UserID: "U001", ServiceType: ""}
		
		err := orderService.CreateOrder(orderKosong)
		
		if err == nil {
			t.Errorf("Harusnya error karena kosong, tapi malah sukses")
		}
	})

	t.Run("Sukses bikin order", func(t *testing.T) {
		orderBenar := &domain.Order{UserID: "U001", ServiceType: "ride"}
		
		mockRepo.EXPECT().Save(orderBenar).Return(nil).Times(1)

		err := orderService.CreateOrder(orderBenar)
		
		if err != nil {
			t.Errorf("Harusnya sukses, tapi malah error: %v", err)
		}
	})
}