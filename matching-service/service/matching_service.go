package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"matching-service/domain"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// Kita hanya perlu define interface KafkaProducer di sini,
// karena MatchingService dan MatchingRepository sudah ada di domain.
type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type matchingService struct {
	repo          domain.MatchingRepository // Panggil dari domain
	kafkaProducer KafkaProducer
	kafkaTopic    string
}

// Kembalikan tipe domain.MatchingService
func NewMatchingService(repo domain.MatchingRepository, producer KafkaProducer, topic string) domain.MatchingService {
	return &matchingService{
		repo:          repo,
		kafkaProducer: producer,
		kafkaTopic:    topic,
	}
}

// MatchOrder mencari driver terdekat secara asynchronous
func (s *matchingService) MatchOrder(ctx context.Context, req *domain.MatchRequest) (*domain.MatchedDriver, error) {
	radius := req.Radius
	if radius == 0 {
		radius = 5000 // default 5 km
	}

	// Cari driver terdekat di DB
	driver, err := s.repo.FindNearestDriver(ctx, req.Latitude, req.Longitude, radius)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tidak ada driver yang ditemukan di sekitar lokasi")
		}
		return nil, fmt.Errorf("gagal mencari driver: %w", err)
	}

	// Goroutine: Kirim pesan Kafka sebagai background process berkala
	go func(bgCtx context.Context, matched *domain.MatchedDriver, orderID uint) {
		eventPayload := map[string]interface{}{
			"event_type": "MATCH_FOUND",
			"order_id":   orderID,
			"driver_id":  matched.DriverID,
			"distance":   matched.Distance,
		}

		eventBytes, _ := json.Marshal(eventPayload)
		err := s.kafkaProducer.WriteMessages(bgCtx, kafka.Message{
			Topic: s.kafkaTopic,
			Key:   []byte(fmt.Sprintf("%d", orderID)),
			Value: eventBytes,
		})

		if err != nil {
			log.Printf("[Error] Gagal publish event MATCH_FOUND untuk order %d: %v", orderID, err)
		} else {
			log.Printf("[Info] Berhasil publish event MATCH_FOUND untuk order %d ke Kafka", orderID)
		}
	}(context.Background(), driver, req.OrderID)

	return driver, nil
}