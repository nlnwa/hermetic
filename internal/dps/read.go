package dps

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func ReadMessages(ctx context.Context, reader *kafka.Reader, callback func(*KafkaResponse) error) error {
	for {
		// TODO log debug
		fmt.Println("Reading next message...")
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var dpsResponse DigitalPreservationSystemResponse

		err = json.Unmarshal(message.Value, &dpsResponse)
		if err != nil {
			return fmt.Errorf("could not unmarshal message at offset '%d': %w", message.Offset, err)
		}

		if !IsWebArchiveOwned(&dpsResponse) {
			// TODO log debug message
			fmt.Printf("Message at offset '%d' is not owned by web archive, skipping offset\n", message.Offset)
			continue
		}

		response := &KafkaResponse{
			Offset:      message.Offset,
			Key:         string(message.Key),
			DPSResponse: dpsResponse,
		}

		err = callback(response)
		if err != nil {
			return err
		}
	}
}
