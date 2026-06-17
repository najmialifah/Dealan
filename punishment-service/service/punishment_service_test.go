package service_test

import (
	"context"
	"testing"

	"github.com/najmialifah/Dealan/punishment-service/domain"
	"github.com/najmialifah/Dealan/punishment-service/mocks"
	"github.com/najmialifah/Dealan/punishment-service/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPunishmentService_ApplyPunishment_Suspended(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPunishmentRepository(ctrl)
	svc := service.NewPunishmentService(mockRepo)

	req := domain.PunishmentRequest{
		AccountID:  "DRIVER-99",
		ReasonCode: "CANCEL_ABUSE",
		Duration:   24, // 24 jam sanksi (Suspended)
	}

	mockRepo.EXPECT().
		StoreViolation(gomock.Any(), req).
		Return("violation-uuid-123", nil).
		Times(1)

	mockRepo.EXPECT().
		UpdateAccountStatus(gomock.Any(), "DRIVER-99", "Suspended").
		Return(nil).
		Times(1)

	res, err := svc.ApplyPunishment(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "violation-uuid-123", res.PenaltyID)
	assert.Equal(t, "Suspended", res.NewAccountStatus)
	assert.Equal(t, "Success", res.Status)
}

func TestPunishmentService_ApplyPunishment_Banned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPunishmentRepository(ctrl)
	svc := service.NewPunishmentService(mockRepo)

	req := domain.PunishmentRequest{
		AccountID:  "USER-88",
		ReasonCode: "FRAUD",
		Duration:   -1, // Banned permanen
	}

	mockRepo.EXPECT().
		StoreViolation(gomock.Any(), req).
		Return("violation-uuid-456", nil).
		Times(1)

	mockRepo.EXPECT().
		UpdateAccountStatus(gomock.Any(), "USER-88", "Banned").
		Return(nil).
		Times(1)

	res, err := svc.ApplyPunishment(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "violation-uuid-456", res.PenaltyID)
	assert.Equal(t, "Banned", res.NewAccountStatus)
	assert.Equal(t, "Success", res.Status)
}