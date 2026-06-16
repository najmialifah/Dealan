package service

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/models"
	"order-service/repository"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer mendefinisikan kontrak untuk mengirim pesan ke Kafka
type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

// OrderService mendefinisikan logika bisnis terkait order
type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}

type orderService struct {
	repo          repository.OrderRepository
	kafkaProducer KafkaProducer
	kafkaTopic    string
}

// NewOrderService menginisialisasi order service
func NewOrderService(repo repository.OrderRepository, kafkaProducer KafkaProducer, kafkaTopic string) OrderService {
	return &orderService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		kafkaTopic:    kafkaTopic,
	}
}

// CreateOrder membuat order baru dan mengirimkan event ORDER_CREATED ke Kafka
func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Tetapkan status default
	order.Status = "PENDING"

	// Simpan ke database melalui repository
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return err
	}

	// Persiapkan event payload
	eventPayload := map[string]interface{}{
		"event_type": "ORDER_CREATED",
		"order_id":   order.ID,
		"user_id":    order.UserID,
		"status":     order.Status,
		"timestamp":  order.CreatedAt,
	}

	eventBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("gagal marshalling event payload: %w", err)
	}

	// Kirim pesan ke Kafka
	err = s.kafkaProducer.WriteMessages(ctx, kafka.Message{
		Topic: s.kafkaTopic,
		Key:   []byte(fmt.Sprintf("%d", order.ID)),
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("gagal publish ke kafka: %w", err)
	}

	return nil
}