package confirmImplmentation

import (
	"context"
	"fmt"
	"hermetic/internal/dps"

	"github.com/segmentio/kafka-go"
)

func ReadConfirmTopic(ctx context.Context, reader *kafka.Reader) error {
	for {
		response, err := dps.ReadMessage(ctx, reader)
		if err != nil {
			return fmt.Errorf("failed to read confirm-topic: `%w`", err)
		}
		
		fmt.Printf("SIP successfully preserved: %+v\n", response)
	}
}
