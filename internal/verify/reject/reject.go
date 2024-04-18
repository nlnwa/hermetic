package rejectImplementation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
	"hermetic/internal/teams"
)

func ReadRejectTopic(ctx context.Context, reader *kafka.Reader, teamsWebhookNotificationUrl string) error {
	err := dps.ReadMessages(ctx, reader, func(response *dps.KafkaResponse) error {
		if err := ProcessMessagesFromRejectTopic(reader, response, teamsWebhookNotificationUrl); err != nil {
			return fmt.Errorf("failed to process message from reject-topic: `%w`", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read reject-topic: `%w`", err)
	}
	return nil
}

func ProcessMessagesFromRejectTopic(reader *kafka.Reader, response *dps.KafkaResponse, teamsWebhookNotificationUrl string) error {
	if response == nil {
		panic("No response found")
	}

	// TODO log more info from reject message in case sending to teams fails
	fmt.Printf("Processing message with ContentCategory: '%s'\n", response.DPSResponse.ContentCategory)

	payload := createTeamsDigitalPreservationSystemFailureMessage(response, reader.Config().Topic, reader.Config().Brokers)
	err := teams.SendMessage(payload, teamsWebhookNotificationUrl)
	if err != nil {
		return fmt.Errorf("failed to send message to teams, cause: `%w`", err)
	}
	return nil
}
