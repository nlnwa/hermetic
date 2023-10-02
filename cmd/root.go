package cmd

import (
	"hermetic/cmd/send"
	"hermetic/cmd/verify"
	"hermetic/internal/common_flags"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "hermetic",
		Short: "hermetic - sends and verifies data for digital storage",
	}
	rootCommand.PersistentFlags().StringSliceVar(&common_flags.KafkaEndpoints, "kafka-endpoints", []string{}, "list of kafka endpoints")
	rootCommand.PersistentFlags().StringVar(&common_flags.TeamsWebhookNotificationUrl, "teams-webhook-notification-url", "", "url to teams webhook for notifications")
	rootCommand.AddCommand(send.NewCommand())
	rootCommand.AddCommand(verify.NewCommand())
	return rootCommand
}
