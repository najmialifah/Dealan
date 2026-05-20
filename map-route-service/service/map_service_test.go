package service

import (
	"context"
	"testing"

	"map-route-service/domain"
	"map-route-service/mocks"

	"go.uber.org/mock/gomock"
)

func TestGetRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMapRepository(ctrl)

	mockRepo.EXPECT().
		GetRoute(gomock.Any(), "A", "B").
		Return(&domain.RouteResponse{
			Distance: 10,
			Duration: 15,
		}, nil)

	svc := NewMapService(mockRepo)

	res, err := svc.GetRoute(context.Background(), domain.RouteRequest{
		Origin:      "A",
		Destination: "B",
	})

	if err != nil {
		t.Errorf("error")
	}

	if res.Distance != 10 {
		t.Errorf("wrong result")
	}
}