package service

import (
	"context"
	"order-service/models"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository adalah mock untuk OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	// Simulasi penyimpanan dengan memberikan ID jika sukses
	if args.Error(0) == nil {
		order.ID = 1
	}
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockKafkaProducer adalah mock untuk KafkaProducer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	// Untuk kemudahan mock, kita anggap hanya mengirim 1 pesan pada satu waktu di test ini
	args := m.Called(ctx, msgs[0])
	return args.Error(0)
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockKafka := new(MockKafkaProducer)
	topic := "ORDER_CREATED_TEST"

	service := NewOrderService(mockRepo, mockKafka, topic)

	order := &models.Order{
		UserID: 123,
		DetailPaket: map[string]interface{}{
			"barang": "Laptop",
			"berat":  2.5,
		},
	}

	// Ekspektasi: repo.CreateOrder dipanggil sekali
	mockRepo.On("CreateOrder", mock.Anything, order).Return(nil)

	// Ekspektasi: kafkaProducer.WriteMessages dipanggil sekali
	mockKafka.On("WriteMessages", mock.Anything, mock.AnythingOfType("kafka.Message")).Return(nil)

	err := service.CreateOrder(context.Background(), order)

	assert.NoError(t, err)
	assert.Equal(t, "PENDING", order.Status)
	assert.Equal(t, uint(1), order.ID)

	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}