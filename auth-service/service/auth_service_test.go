package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
	"github.com/shakilaaulia/Dealan/auth-service/mocks"
	"github.com/shakilaaulia/Dealan/auth-service/service"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login_RepositoryError(t *testing.T) {
	// 1. Setup GoMock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Setup Mock Repository
	mockRepo := mocks.NewMockAuthRepository(ctrl)

	// 3. Setup Service with Mock Repository (Dependency Injection)
	authSvc := service.NewAuthService(mockRepo)

	// 4. Define Scenario: Mock behavior to return error on GetUserByEmail
	expectedErr := errors.New("database connection failed")
	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(nil, expectedErr).
		Times(1)

	// 5. Execute Action
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := authSvc.Login(context.Background(), req)

	// 6. Assert Expectations
	if err == nil {
		t.Errorf("expected error but got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Errorf("expected response nil, got %v", resp)
	}
}
