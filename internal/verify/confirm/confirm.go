package confirmImplmentation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	kafkaHelpers "hermetic/internal/kafka"
	"time"
)

func isWebArchiveOwned(message kafkaResponse) bool {
	return message.DPSResponse.ContentCategory == "nettarkiv"
}

func Verify(confirmTopicName string, kafkaEndpoints []string) error {
	reader := kafkaHelpers.MessageReader{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: kafkaEndpoints,
			Topic:   confirmTopicName,
			GroupID: "nettarkivet-hermetic-verify",
		}),
	}
	readTimeout := 10 * time.Second
	cycleSleepDuration := 1 * time.Minute

	for {
		message, err := reader.ReadMessageWithTimeout(readTimeout)
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Reading message timed out, sleeping for '%s'\n", cycleSleepDuration)
			time.Sleep(cycleSleepDuration)
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to read message, cause: `%w`", err)
		}
		var parsedMessage digitalPreservationSystemResponse
		err = json.Unmarshal(message.Value, &parsedMessage)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				fmt.Printf("Could not read message at offset '%d', syntax error in message, skipping offset\n", message.Offset)
				continue
			}
			return fmt.Errorf("failed to unmarshal json, original error: '%w'", err)
		}
		response := kafkaResponse{
			Offset:      message.Offset,
			Key:         string(message.Key),
			DPSResponse: parsedMessage,
		}
		if !isWebArchiveOwned(response) {
			fmt.Printf("Skipping message with ContentCategory: '%s'\n", response.DPSResponse.ContentCategory)
			continue
		}
	}
}
