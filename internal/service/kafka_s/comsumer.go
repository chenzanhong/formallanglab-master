package kafka

import "context"

type KafkaEmailConsumerService interface {
	ReadEmail(ctx context.Context) (*KafkaEmailEvent, error)
}
