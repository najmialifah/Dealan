package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/shakilaaulia/Dealan/auth-service/service"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer membuat instance baru EventProducer untuk Apache Kafka
func NewKafkaProducer(brokers []string) service.EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// PublishUserCreated mengirimkan event USER_CREATED ke topik Kafka "user.created"
func (p *kafkaProducer) PublishUserCreated(ctx context.Context, id, name, email, phone string) error {
	msgData := map[string]interface{}{
		"event": "USER_CREATED",
		"payload": map[string]interface{}{
			"id":       id,
			"nama":     name,
			"email":    email,
			"nomor_hp": phone,
		},
	}

	bytes, err := json.Marshal(msgData)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user.created",
		Key:   []byte(id),
		Value: bytes,
	})
}

// PublishDriverCreated mengirimkan event DRIVER_CREATED ke topik Kafka "driver.created"
func (p *kafkaProducer) PublishDriverCreated(ctx context.Context, id, name, email, phone string) error {
	msgData := map[string]interface{}{
		"event": "DRIVER_CREATED",
		"payload": map[string]interface{}{
			"id":       id,
			"nama":     name,
			"email":    email,
			"nomor_hp": phone,
		},
	}

	bytes, err := json.Marshal(msgData)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: "driver.created",
		Key:   []byte(id),
		Value: bytes,
	})
}
