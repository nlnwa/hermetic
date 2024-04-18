package confirmImplmentation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
)

func ReadConfirmTopic(ctx context.Context, reader *kafka.Reader) error {
	err := dps.ReadMessages(ctx, reader, ProcessMessagesFromConfirmTopic)
	if err != nil {
		return fmt.Errorf("failed to read confirm-topic: `%w`", err)
	}
	return nil
}

func ProcessMessagesFromConfirmTopic(response *dps.KafkaResponse) error {
	if response == nil {
		panic("No response found")
	}

	fmt.Printf("SIP successfully preserved: %+v\n", response)

	return nil
}
