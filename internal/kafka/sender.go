package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Sender struct {
	Writer *kafka.Writer
}

func (sender *Sender) SendMessageToKafkaTopic(payload []byte) error {
	kafkaMessageUuid := uuid.New()
	slog.Info("Sending message to Kafka topic", "UUID", kafkaMessageUuid)
	kafkaMessageUuidBytes, err := kafkaMessageUuid.MarshalText()
	if err != nil {
		return fmt.Errorf("failed to marshal uuid, original error: '%w'", err)
	}

	err = sender.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   kafkaMessageUuidBytes,
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write messages, original error: '%w'", err)
	}

	return nil
}
