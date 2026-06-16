package service

import (
	"context"
	"matching-service/domain" // Sudah diganti ke domain
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockMatchingRepository
type MockMatchingRepository struct {
	mock.Mock
}

// radiusMeters diubah jadi int agar sesuai dengan kontrak di domain
func (m *MockMatchingRepository) FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters int) (*domain.MatchedDriver, error) {
	args := m.Called(ctx, lat, lon, radiusMeters)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.MatchedDriver), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockKafkaProducer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs[0])
	return args.Error(0)
}

func TestMatchingService_MatchOrder_Success(t *testing.T) {
	mockRepo := new(MockMatchingRepository)
	mockKafka := new(MockKafkaProducer)
	svc := NewMatchingService(mockRepo, mockKafka, "MATCH_TEST")

	req := &domain.MatchRequest{
		OrderID:   10,
		Latitude:  -6.2,
		Longitude: 106.8,
		Radius:    2000,
	}

	expectedDriver := &domain.MatchedDriver{
		DriverID:  "driver-100", // Diubah jadi string (UUID)
		Latitude:  -6.201,
		Longitude: 106.801,
		Distance:  150,
	}

	// 2000 diubah jadi int (tidak perlu float64)
	mockRepo.On("FindNearestDriver", mock.Anything, -6.2, 106.8, 2000).Return(expectedDriver, nil)
	mockKafka.On("WriteMessages", mock.Anything, mock.AnythingOfType("kafka.Message")).Return(nil)

	driver, err := svc.MatchOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, "driver-100", driver.DriverID) // Diubah ke string

	time.Sleep(10 * time.Millisecond) // wait goroutine to fire
	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestMatchingService_MatchOrder_NotFound(t *testing.T) {
	mockRepo := new(MockMatchingRepository)
	mockKafka := new(MockKafkaProducer)
	svc := NewMatchingService(mockRepo, mockKafka, "MATCH_TEST")

	req := &domain.MatchRequest{
		OrderID:   10,
		Latitude:  -6.2,
		Longitude: 106.8,
	}

	// 5000 sebagai int
	mockRepo.On("FindNearestDriver", mock.Anything, -6.2, 106.8, 5000).Return((*domain.MatchedDriver)(nil), gorm.ErrRecordNotFound)

	driver, err := svc.MatchOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, driver)
	assert.Equal(t, "tidak ada driver yang ditemukan di sekitar lokasi", err.Error())
	mockRepo.AssertExpectations(t)
	// Kafka shouldn't be called
	mockKafka.AssertNotCalled(t, "WriteMessages")
}