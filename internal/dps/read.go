package dps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

func ReadMessages(ctx context.Context, reader *kafka.Reader) (*KafkaResponse, error) {
	for {
		slog.Info("Reading next message...")
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var dpsResponse DigitalPreservationSystemResponse

		err = json.Unmarshal(message.Value, &dpsResponse)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				slog.Error("Could not read message at offset, syntax error in message, skipping offset", "offset", message.Offset)
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal json, original error: '%w'", err)
		}

		if !IsWebArchiveOwned(&dpsResponse) {
			slog.Info("Message at offset is not owned by web archive, skipping offset", "offset", message.Offset)
			continue
		}

		return &KafkaResponse{
			Offset:      message.Offset,
			Key:         string(message.Key),
			DPSResponse: dpsResponse,
		}, nil
	}
}
