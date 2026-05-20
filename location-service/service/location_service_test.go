package service

import (
	"location-service/domain"
	"location-service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestUpdateDriverLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockLocationRepository(ctrl)
	locService := LocationService{Repo: mockRepo}

	t.Run("Gagal Latitude Ngaco", func(t *testing.T) {
		loc := &domain.LocationUpdate{EntityID: "D001", Latitude: 999.0, Longitude: 100.0}
		err := locService.UpdateDriverLocation(loc)
		if err == nil {
			t.Errorf("Harusnya error karena latitude 999 tidak ada di peta")
		}
	})

	t.Run("Sukses Update Lokasi", func(t *testing.T) {
		loc := &domain.LocationUpdate{EntityID: "D001", Latitude: -6.200000, Longitude: 106.816666} // Koordinat normal
		
		mockRepo.EXPECT().SaveLocation(loc).Return(nil).Times(1)

		err := locService.UpdateDriverLocation(loc)
		if err != nil {
			t.Errorf("Harusnya sukses update lokasi, tapi malah error")
		}
	})
}