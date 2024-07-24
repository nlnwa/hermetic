package reject

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"hermetic/internal/common_flags"
	"hermetic/internal/teams"
	rejectImplementation "hermetic/internal/verify/reject"
	"os"
	"os/signal"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "reject",
		Short: "Continuously report all rejected data",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndReadRejectTopic,
	}
	rejectTopicFlagName := "reject-topic"
	command.Flags().String(rejectTopicFlagName, "", "name of reject-topic")
	if err := command.MarkFlagRequired(rejectTopicFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", rejectTopicFlagName, err))
	}
	return command
}

func parseArgumentsAndReadRejectTopic(cmd *cobra.Command, args []string) error {
	rejectTopicName, err := cmd.Flags().GetString("reject-topic")
	if err != nil {
		return fmt.Errorf("failed to get reject-topic flag, cause: `%w`", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: common_flags.KafkaEndpoints,
		Topic:   rejectTopicName,
		GroupID: "nettarkivet-hermetic-verify-reject",
	})

	err = rejectImplementation.ReadRejectTopic(ctx, reader, common_flags.TeamsWebhookNotificationUrl)
	if err != nil {
		err = fmt.Errorf("verification error, cause: `%w`", err)
		fmt.Printf("Sending error message to Teams\n")
		teamsErrorMessage := teams.CreateGeneralFailureMessage(err)
		if err := teams.SendMessage(teamsErrorMessage, common_flags.TeamsWebhookNotificationUrl); err != nil {
			err = fmt.Errorf("failed to send error message to Teams, cause: `%w`", err)
			fmt.Printf("%s\n", err)
		}
		return err
	}
	return err
}
