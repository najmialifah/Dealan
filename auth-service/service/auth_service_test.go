package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
	"github.com/shakilaaulia/Dealan/auth-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockAuthRepository adalah implementasi mock untuk repository.AuthRepository menggunakan testify/mock
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetCredentialByEmail(ctx context.Context, email string) (*domain.AuthCredential, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthCredential), args.Error(1)
}

func (m *MockAuthRepository) CreateCredential(ctx context.Context, cred *domain.AuthCredential) error {
	args := m.Called(ctx, cred)
	return args.Error(0)
}

func (m *MockAuthRepository) CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	args := m.Called(ctx, rt)
	return args.Error(0)
}

func (m *MockAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

// MockEventProducer adalah implementasi mock untuk service.EventProducer menggunakan testify/mock
type MockEventProducer struct {
	mock.Mock
}

func (m *MockEventProducer) PublishUserCreated(ctx context.Context, id, name, email, phone string) error {
	args := m.Called(ctx, id, name, email, phone)
	return args.Error(0)
}

func (m *MockEventProducer) PublishDriverCreated(ctx context.Context, id, name, email, phone string) error {
	args := m.Called(ctx, id, name, email, phone)
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockProducer := new(MockEventProducer)
	jwtSecret := "test_secret_12345!"
	jwtDuration := 1 * time.Hour

	authSvc := service.NewAuthService(mockRepo, mockProducer, jwtSecret, jwtDuration)

	req := domain.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Role:     domain.RoleUser,
		Nama:     "Budi Santoso",
		NomorHP:  "0812345678",
	}

	// Mock behavior: GetCredentialByEmail mengembalikan nil (belum terdaftar)
	mockRepo.On("GetCredentialByEmail", mock.Anything, req.Email).Return(nil, errors.New("record not found"))
	// Mock behavior: CreateCredential sukses
	mockRepo.On("CreateCredential", mock.Anything, mock.Anything).Return(nil)
	// Mock behavior: PublishUserCreated sukses
	mockProducer.On("PublishUserCreated", mock.Anything, mock.Anything, req.Nama, req.Email, req.NomorHP).Return(nil)

	resp, err := authSvc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, string(domain.RoleUser), resp.Role)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockProducer := new(MockEventProducer)
	jwtSecret := "test_secret_12345!"

	authSvc := service.NewAuthService(mockRepo, mockProducer, jwtSecret, 1*time.Hour)

	req := domain.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Role:     domain.RoleUser,
	}

	existingCred := &domain.AuthCredential{
		Email: "existing@example.com",
	}

	// Mock behavior: GetCredentialByEmail menemukan email tersebut
	mockRepo.On("GetCredentialByEmail", mock.Anything, req.Email).Return(existingCred, nil)

	resp, err := authSvc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email sudah terdaftar", err.Error())
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockProducer := new(MockEventProducer)
	jwtSecret := "test_secret_12345!"

	authSvc := service.NewAuthService(mockRepo, mockProducer, jwtSecret, 1*time.Hour)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretPassword"), bcrypt.DefaultCost)

	cred := &domain.AuthCredential{
		AccountID:    "acc-123",
		Email:        "user@example.com",
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleUser,
		IsActive:     true,
	}

	req := domain.LoginRequest{
		Email:    "user@example.com",
		Password: "secretPassword",
	}

	mockRepo.On("GetCredentialByEmail", mock.Anything, req.Email).Return(cred, nil)

	resp, err := authSvc.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "acc-123", resp.AccountID)
	assert.Equal(t, string(domain.RoleUser), resp.Role)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	mockProducer := new(MockEventProducer)
	jwtSecret := "test_secret_12345!"

	authSvc := service.NewAuthService(mockRepo, mockProducer, jwtSecret, 1*time.Hour)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretPassword"), bcrypt.DefaultCost)

	cred := &domain.AuthCredential{
		Email:        "user@example.com",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	req := domain.LoginRequest{
		Email:    "user@example.com",
		Password: "wrongPassword",
	}

	mockRepo.On("GetCredentialByEmail", mock.Anything, req.Email).Return(cred, nil)

	resp, err := authSvc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "kredensial tidak valid", err.Error())
}
