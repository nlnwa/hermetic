package confirm

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"hermetic/internal/common_flags"
	"hermetic/internal/teams"
	confirmImplementation "hermetic/internal/verify/confirm"
	"os"
	"os/signal"
)

const (
	confirmTopicFlagName           string = "confirm-topic"
	confirmMessageReceiverFlagName string = "confirm-message-receiver"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "confirm",
		Short: "Continuously report all successfully preserved data",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndReadConfirmTopic,
	}
	command.Flags().String(confirmTopicFlagName, "", "name of confirm-topic")
	if err := command.MarkFlagRequired(confirmTopicFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", confirmTopicFlagName, err))
	}

	command.Flags().String(confirmMessageReceiverFlagName, "", "optional URL for confirm message receiver")

	return command
}

func parseArgumentsAndReadConfirmTopic(cmd *cobra.Command, args []string) error {
	confirmTopicName, err := cmd.Flags().GetString(confirmTopicFlagName)
	if err != nil {
		return fmt.Errorf("failed to get confirm-topic flag, cause: `%w`", err)
	}

	confirmMessageReceiverUrl, err := cmd.Flags().GetString(confirmMessageReceiverFlagName)
	if err != nil {
		return fmt.Errorf("failed to get confirm-message-receiver flag, cause: `%w`", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: common_flags.KafkaEndpoints,
		Topic:   confirmTopicName,
		GroupID: "nettarkivet-hermetic-verify-confirm",
	})

	if confirmMessageReceiverUrl != "" {
		err = confirmImplementation.ReadConfirmTopic(ctx, reader, confirmMessageReceiverUrl)
	} else {
		err = confirmImplementation.ReadConfirmTopic(ctx, reader)
	}
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
