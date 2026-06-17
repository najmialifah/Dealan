package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/najmialifah/Dealan/user-service/service"
	"github.com/segmentio/kafka-go"
)

// UserCreatedEvent merepresentasikan payload event Kafka yang diterima
type UserCreatedEvent struct {
	Event   string `json:"event"`
	Payload struct {
		ID      string `json:"id"`
		Nama    string `json:"nama"`
		Email   string `json:"email"`
		NomorHP string `json:"nomor_hp"`
	} `json:"payload"`
}

// KafkaConsumer mendengarkan event register dari Kafka untuk user-service
type KafkaConsumer struct {
	reader *kafka.Reader
	svc    service.UserService
}

// NewKafkaConsumer membuat instance baru dari KafkaConsumer
func NewKafkaConsumer(brokers []string, groupID string, svc service.UserService) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    "user.created",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		svc: svc,
	}
}

// Start mulai mendengarkan event Kafka secara asinkron
func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Memulai Kafka Consumer untuk topic: user.created...")
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

			var event UserCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("Gagal melakukan unmarshal payload event:", err)
				continue
			}

			if event.Event == "USER_CREATED" {
				log.Printf("Menerima event USER_CREATED untuk ID: %s, Nama: %s", event.Payload.ID, event.Payload.Nama)
				err := c.svc.CreateUser(ctx, event.Payload.ID, event.Payload.Nama, event.Payload.Email, event.Payload.NomorHP)
				if err != nil {
					log.Printf("Gagal membuat user dari event Kafka: %v", err)
				}
			}
		}
	}
}
