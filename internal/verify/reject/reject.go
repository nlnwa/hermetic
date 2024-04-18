package rejectImplementation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
	"hermetic/internal/teams"
)

func ReadRejectTopic(ctx context.Context, reader *kafka.Reader, teamsWebhookNotificationUrl string) error {
	err := dps.ReadMessages(ctx, reader, func(response *dps.KafkaResponse, err error) {
		ProcessMessagesFromRejectTopic(reader, response, teamsWebhookNotificationUrl, err)
	})
	if err != nil {
		return fmt.Errorf("failed to read reject-topic: `%w`", err)
	}
	return nil
}
func ProcessMessagesFromRejectTopic(reader *kafka.Reader, response *dps.KafkaResponse, teamsWebhookNotificationUrl string, err error) {
	if err != nil {
		fmt.Errorf("failed to read message, cause: `%w`", err)
	}
	if response == nil {
		fmt.Println("No response found")
	} else {
		fmt.Printf("Processing message with ContentCategory: '%s'\n", response.DPSResponse.ContentCategory)
		payload := createTeamsDigitalPreservationSystemFailureMessage(response, reader.Config().Topic, reader.Config().Brokers)
		err = teams.SendMessage(payload, teamsWebhookNotificationUrl)
		if err != nil {
			fmt.Errorf("failed to send message to teams, cause: `%w`", err)
		}
	}
}
