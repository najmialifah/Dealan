package service

import (
	"context"
	"testing"

	"github.com/shakilaaulia/Dealan/punishment-service/domain"
	"github.com/shakilaaulia/Dealan/punishment-service/mocks"
	"go.uber.org/mock/gomock"
)

func TestPunishmentService_Apply_ExpectedFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPunishmentRepository(ctrl)
	svc := NewPunishmentService(mockRepo)

	req := domain.PunishmentRequest{
		AccountID:  "DRIVER-99",
		ReasonCode: "LATE_CANCEL",
		Duration:   24,
	}

	_, err := svc.ApplyPunishment(context.Background(), req)

	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Expected 'not implemented' error, got %v", err)
	}
}