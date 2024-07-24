package dps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func ReadMessages(ctx context.Context, reader *kafka.Reader) (*KafkaResponse, error) {
	for {
		fmt.Println("Reading next message...")
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var dpsResponse DigitalPreservationSystemResponse

		err = json.Unmarshal(message.Value, &dpsResponse)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				fmt.Printf("Could not read message at offset '%d', syntax error in message, skipping offset\n", message.Offset)
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal json, original error: '%w'", err)
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
