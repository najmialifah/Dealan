package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"matching-service/models"
	"matching-service/repository"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type MatchingService interface {
	MatchOrder(ctx context.Context, req *models.MatchRequest) (*models.MatchedDriver, error)
}

type matchingService struct {
	repo          repository.MatchingRepository
	kafkaProducer KafkaProducer
	kafkaTopic    string
}

func NewMatchingService(repo repository.MatchingRepository, producer KafkaProducer, topic string) MatchingService {
	return &matchingService{
		repo:          repo,
		kafkaProducer: producer,
		kafkaTopic:    topic,
	}
}

// MatchOrder mencari driver terdekat secara asynchronous atau synchronous dan mengirim event MATCH_FOUND
func (s *matchingService) MatchOrder(ctx context.Context, req *models.MatchRequest) (*models.MatchedDriver, error) {
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

	// Goroutine: Kirim pesan Kafka sebagai background process berkala / event driven agar tidak memblokir API HTTP.
	// Penggunaan goroutine sangat disarankan jika proses notifikasi butuh waktu/retries.
	go func(bgCtx context.Context, matched *models.MatchedDriver, orderID uint) {
		eventPayload := map[string]interface{}{
			"event_type":  "MATCH_FOUND",
			"order_id":    orderID,
			"driver_id":   matched.DriverID,
			"distance":    matched.Distance,
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