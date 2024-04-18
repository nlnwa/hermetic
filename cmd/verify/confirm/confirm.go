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

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "confirm",
		Short: "Continuously report all successfully preserved data",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallVerify,
	}
	confirmTopicFlagName := "confirm-topic"
	command.Flags().String(confirmTopicFlagName, "", "name of confirm-topic")
	if err := command.MarkFlagRequired(confirmTopicFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", confirmTopicFlagName, err))
	}
	return command
}

func parseArgumentsAndCallVerify(cmd *cobra.Command, args []string) error {
	confirmTopicName, err := cmd.Flags().GetString("confirm-topic")
	if err != nil {
		return fmt.Errorf("failed to get confirm-topic flag, cause: `%w`", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: common_flags.KafkaEndpoints,
		Topic:   confirmTopicName,
		GroupID: "nettarkivet-hermetic-verify-confirm",
	})

	err = confirmImplementation.ReadConfirmTopic(ctx, reader)
	if err != nil {
		fmt.Printf("Verification error: %v\n", err)
		teamsErrorMessage := teams.CreateGeneralFailureMessage(err)
		if err := teams.SendMessage(teamsErrorMessage, common_flags.TeamsWebhookNotificationUrl); err != nil {
			fmt.Printf("Failed to send error message to Teams: %v\n", err)
		}
	}
	return err
}
