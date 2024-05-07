package dps

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func ReadMessage(ctx context.Context, reader *kafka.Reader) (*KafkaResponse, error) {
	for {
		fmt.Println("Reading next message...")
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var dpsResponse DigitalPreservationSystemResponse

		err = json.Unmarshal(message.Value, &dpsResponse)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal message at offset '%d': %w", message.Offset, err)
		}

		if !IsWebArchiveOwned(&dpsResponse) {
			fmt.Printf("Message at offset '%d' is not owned by web archive, skipping offset\n", message.Offset)
			continue
		}

		return &KafkaResponse{
			Offset:      message.Offset,
			Key:         string(message.Key),
			DPSResponse: dpsResponse,
		}, nil
	}
}
