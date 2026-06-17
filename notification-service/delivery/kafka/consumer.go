package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/najmialifah/Dealan/notification-service/domain"
	"github.com/najmialifah/Dealan/notification-service/service"
	kafkaGo "github.com/segmentio/kafka-go"
)

// PaymentCompletedEvent mencerminkan event pembayaran selesai dari broker Kafka
type PaymentCompletedEvent struct {
	Event     string `json:"event"`
	Payload   struct {
		TransactionID  string  `json:"transaction_id"`
		OrderID        string  `json:"order_id"`
		Status         string  `json:"status"`
		DriverEarnings float64 `json:"driver_earnings"`
	} `json:"payload"`
}

type KafkaConsumer struct {
	reader *kafkaGo.Reader
	svc    service.NotificationService
}

// NewKafkaConsumer menginisialisasi Kafka Consumer baru
func NewKafkaConsumer(brokers []string, topic string, groupID string, svc service.NotificationService) *KafkaConsumer {
	reader := kafkaGo.NewReader(kafkaGo.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader: reader,
		svc:    svc,
	}
}

// Start mulai mendengarkan event dari Kafka secara terus-menerus
func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Printf("[Kafka Consumer] Mendengarkan topik: %s\n", c.reader.Config().Topic)
	for {
		select {
		case <-ctx.Done():
			log.Println("[Kafka Consumer] Dihentikan.")
			c.reader.Close()
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("[Kafka Consumer] Gagal membaca pesan: %v\n", err)
				continue
			}

			log.Printf("[Kafka Consumer] Pesan diterima [offset %d]: %s\n", m.Offset, string(m.Value))

			var event PaymentCompletedEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("[Kafka Consumer] Gagal unmarshal event: %v\n", err)
				continue
			}

			// Trigger notifikasi jika status pembayaran sukses
			if event.Event == "PAYMENT_COMPLETED" && event.Payload.Status == "success" {
				req := domain.NotificationRequest{
					TargetID:   "user-placeholder", // Dalam skenario nyata, ID ini dicari dari database order
					Title:      "Pembayaran Sukses!",
					Body:       fmt.Sprintf("Pembayaran untuk Order %s senilai Rp %v berhasil.", event.Payload.OrderID, event.Payload.DriverEarnings),
					ActionLink: fmt.Sprintf("/orders/%s", event.Payload.OrderID),
				}

				_, err := c.svc.SendNotification(ctx, req)
				if err != nil {
					log.Printf("[Kafka Consumer] Gagal mengirim notifikasi: %v\n", err)
				} else {
					log.Printf("[Kafka Consumer] Notifikasi terkirim sukses untuk Order %s\n", event.Payload.OrderID)
				}
			}
		}
	}
}
