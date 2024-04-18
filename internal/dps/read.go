package dps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func ReadMessages(ctx context.Context, reader *kafka.Reader, callback func(*KafkaResponse, error)) error {
	for {
		fmt.Println("Reading next message...")
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var dpsResponse DigitalPreservationSystemResponse

		err = json.Unmarshal(message.Value, &dpsResponse)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				fmt.Printf("Could not read message at offset '%d', syntax error in message, skipping offset\n", message.Offset)
				continue
			}
			fmt.Println("failed to unmarshal json: %w", err)
			callback(nil, fmt.Errorf("failed to unmarshal json: '%w'", err))
			continue
		}

		if !IsWebArchiveOwned(&dpsResponse) {
			fmt.Printf("Message at offset '%d' is not owned by web archive, skipping offset\n", message.Offset)
			continue
		}

		response := &KafkaResponse{
			Offset:      message.Offset,
			Key:         string(message.Key),
			DPSResponse: dpsResponse,
		}

		callback(response, nil)
	}
}
