package confirmImplmentation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
)

func ReadConfirmTopic(ctx context.Context, reader *kafka.Reader) error {
	for {
		response, err := dps.ReadMessages(ctx, reader)
		if err != nil {
			return fmt.Errorf("failed to read message from confirm-topic: `%w`", err)
		}
		fmt.Printf("SIP successfully preserved: %+v\n", response)
	}
}
