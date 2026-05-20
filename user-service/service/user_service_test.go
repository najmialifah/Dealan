package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/user-service/domain"
	"github.com/shakilaaulia/Dealan/user-service/mocks"
	"github.com/shakilaaulia/Dealan/user-service/service"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetProfile_RepositoryError(t *testing.T) {
	// 1. Setup GoMock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Setup Mock Repository
	mockRepo := mocks.NewMockUserRepository(ctrl)

	// 3. Setup Service
	userSvc := service.NewUserService(mockRepo)

	// 4. Define Mock Behavior
	expectedErr := errors.New("user not found")
	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), "123").
		Return(nil, expectedErr).
		Times(1)

	// 5. Execute
	profile, err := userSvc.GetProfile(context.Background(), "123")

	// 6. Assert
	if err == nil {
		t.Errorf("expected error but got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if profile != nil {
		t.Errorf("expected profile nil, got %v", profile)
	}
}

func TestUserService_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := service.NewUserService(mockRepo)

	expectedUser := &domain.User{
		ID:    "123",
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), "123").
		Return(expectedUser, nil).
		Times(1)

	profile, err := userSvc.GetProfile(context.Background(), "123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if profile == nil || profile.Name != expectedUser.Name {
		t.Errorf("expected profile name %s, got %v", expectedUser.Name, profile)
	}
}
