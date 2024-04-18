package rejectImplementation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
	"hermetic/internal/teams"
)

func ReadRejectTopic(ctx context.Context, reader *kafka.Reader, teamsWebhookNotificationUrl string) error {
	for {
		response, err := dps.ReadMessages(ctx, reader)
		if err != nil {
			return fmt.Errorf("failed to read message from reject-topic: `%w`", err)
		}
		if err := ProcessMessagesFromRejectTopic(reader, response, teamsWebhookNotificationUrl); err != nil {
			return fmt.Errorf("failed to process message from reject-topic: `%w`", err)
		}
	}
}

func ProcessMessagesFromRejectTopic(reader *kafka.Reader, response *dps.KafkaResponse, teamsWebhookNotificationUrl string) error {
	fmt.Printf("Processing message with ContentCategory: '%s'\n", response.DPSResponse.ContentCategory)
	payload := createTeamsDigitalPreservationSystemFailureMessage(response, reader.Config().Topic, reader.Config().Brokers)
	err := teams.SendMessage(payload, teamsWebhookNotificationUrl)
	if err != nil {
		return fmt.Errorf("failed to send message to teams, cause: `%w`", err)
	}
	return nil
}
