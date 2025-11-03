// backend/internal/service/kafka_s/producer.go
package kafka

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

var (
	producer *kafka.Writer
)

func InitProducerDefault() {
	topic := os.Getenv("KAFKA_TOPIC")
	brokers:= strings.Split(strings.TrimSpace(os.Getenv("KAFKA_BROKERS")),",")
	producer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne, // 根据需要调整
	}
}

func InitProducer(brokers []string, topic string) {
	producer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // 根据需要调整
	}
}

func CloseProducer() {
	if producer != nil {
		producer.Close()
	}
}

func SendEmailEvent(ctx context.Context, event *EmailEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return producer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}
