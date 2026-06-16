package service

import (
	"context"
	"matching-service/models"
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

func (m *MockMatchingRepository) FindNearestDriver(ctx context.Context, lat, lon float64, radiusMeters float64) (*models.MatchedDriver, error) {
	args := m.Called(ctx, lat, lon, radiusMeters)
	if args.Get(0) != nil {
		return args.Get(0).(*models.MatchedDriver), args.Error(1)
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

	req := &models.MatchRequest{
		OrderID:   10,
		Latitude:  -6.2,
		Longitude: 106.8,
		Radius:    2000,
	}

	expectedDriver := &models.MatchedDriver{
		DriverID:  100,
		Latitude:  -6.201,
		Longitude: 106.801,
		Distance:  150,
	}

	mockRepo.On("FindNearestDriver", mock.Anything, -6.2, 106.8, float64(2000)).Return(expectedDriver, nil)
	mockKafka.On("WriteMessages", mock.Anything, mock.AnythingOfType("kafka.Message")).Return(nil)

	driver, err := svc.MatchOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, uint(100), driver.DriverID)

	time.Sleep(10 * time.Millisecond) // wait goroutine to fire
	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestMatchingService_MatchOrder_NotFound(t *testing.T) {
	mockRepo := new(MockMatchingRepository)
	mockKafka := new(MockKafkaProducer)
	svc := NewMatchingService(mockRepo, mockKafka, "MATCH_TEST")

	req := &models.MatchRequest{
		OrderID:   10,
		Latitude:  -6.2,
		Longitude: 106.8,
	}

	// 5000 is default radius when not provided
	mockRepo.On("FindNearestDriver", mock.Anything, -6.2, 106.8, float64(5000)).Return((*models.MatchedDriver)(nil), gorm.ErrRecordNotFound)

	driver, err := svc.MatchOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, driver)
	assert.Equal(t, "tidak ada driver yang ditemukan di sekitar lokasi", err.Error())
	mockRepo.AssertExpectations(t)
	// Kafka shouldn't be called
	mockKafka.AssertNotCalled(t, "WriteMessages")
}