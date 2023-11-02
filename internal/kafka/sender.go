package kafka

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Sender struct {
	Writer *kafka.Writer
}

func (sender *Sender) SendMessageToKafkaTopic(payload []byte) error {
	kafkaMessageUuid := uuid.New()
	fmt.Printf("Sending message with uuid %s\n", kafkaMessageUuid)
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
