package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/najmialifah/Dealan/user-service/domain"
	"github.com/najmialifah/Dealan/user-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah mock untuk repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserService_GetProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	expectedUser := &domain.User{
		ID:      "user-123",
		Nama:    "Budi Santoso",
		Email:   "budi@example.com",
		NomorHP: "0812345678",
	}

	mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(expectedUser, nil)

	profile, err := userSvc.GetProfile(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedUser.Nama, profile.Nama)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	mockRepo.On("GetUserByID", mock.Anything, "invalid-id").Return(nil, errors.New("user not found"))

	profile, err := userSvc.GetProfile(context.Background(), "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Equal(t, "user not found", err.Error())
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	existingUser := &domain.User{
		ID:      "user-123",
		Nama:    "Budi S",
		NomorHP: "0812345678",
	}

	req := domain.UpdateProfileRequest{
		Nama:       "Budi Santoso",
		NomorHP:    "0899999999",
		Alamat:     "Jl. Kebon Jeruk No. 12",
		FotoProfil: "http://storage.com/budi.jpg",
	}

	mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(existingUser, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	err := userSvc.UpdateProfile(context.Background(), "user-123", req)

	assert.NoError(t, err)
	assert.Equal(t, "Budi Santoso", existingUser.Nama)
	assert.Equal(t, "0899999999", existingUser.NomorHP)
	assert.Equal(t, "Jl. Kebon Jeruk No. 12", existingUser.Alamat)
	assert.Equal(t, "http://storage.com/budi.jpg", existingUser.FotoProfil)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	mockRepo.On("GetUserByID", mock.Anything, "new-uuid").Return(nil, errors.New("record not found"))
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == "new-uuid" && user.Nama == "Siti Aminah"
	})).Return(nil)

	err := userSvc.CreateUser(context.Background(), "new-uuid", "Siti Aminah", "siti@example.com", "0811223344")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
