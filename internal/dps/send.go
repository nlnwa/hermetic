package dps

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func CreateUuid() ([]byte, error) {
	id := uuid.New()
	return id.MarshalText()
}

func Send(ctx context.Context, w *kafka.Writer, msg Message) error {
	key, err := CreateUuid()
	if err != nil {
		return err
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   key,
		Value: value,
	}

	return w.WriteMessages(ctx, message)
}
