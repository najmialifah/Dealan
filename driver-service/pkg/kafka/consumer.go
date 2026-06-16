package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/shakilaaulia/Dealan/driver-service/service"
)

// DriverCreatedEvent merepresentasikan payload event Kafka yang diterima
type DriverCreatedEvent struct {
	Event   string `json:"event"`
	Payload struct {
		ID      string `json:"id"`
		Nama    string `json:"nama"`
		Email   string `json:"email"`
		NomorHP string `json:"nomor_hp"`
	} `json:"payload"`
}

// KafkaConsumer mendengarkan event register dari Kafka untuk driver-service
type KafkaConsumer struct {
	reader *kafka.Reader
	svc    service.DriverService
}

// NewKafkaConsumer membuat instance baru dari KafkaConsumer untuk driver-service
func NewKafkaConsumer(brokers []string, groupID string, svc service.DriverService) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    "driver.created",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		svc: svc,
	}
}

// Start mulai mendengarkan event Kafka secara asinkron
func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Memulai Kafka Consumer untuk topic: driver.created...")
	defer c.reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Menghentikan Kafka Consumer...")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Gagal membaca pesan Kafka:", err)
				continue
			}

			var event DriverCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("Gagal melakukan unmarshal payload event:", err)
				continue
			}

			if event.Event == "DRIVER_CREATED" {
				log.Printf("Menerima event DRIVER_CREATED untuk ID: %s, Nama: %s", event.Payload.ID, event.Payload.Nama)
				err := c.svc.CreateDriver(ctx, event.Payload.ID, event.Payload.Nama, event.Payload.Email, event.Payload.NomorHP)
				if err != nil {
					log.Printf("Gagal membuat driver dari event Kafka: %v", err)
				}
			}
		}
	}
}
