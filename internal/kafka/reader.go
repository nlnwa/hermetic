package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageReader struct {
	Reader *kafka.Reader
}

func (reader *MessageReader) ReadMessageWithTimeout(timeout time.Duration) (message kafka.Message, err error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return reader.Reader.ReadMessage(ctxTimeout)
}
