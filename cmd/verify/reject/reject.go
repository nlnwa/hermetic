package reject

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nlnwa/hermetic/cmd/internal/cmdutil"
	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/internal/dps"
	"github.com/nlnwa/hermetic/internal/teams"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command) {
	flags.AddKafkaFlags(cmd)
}

type RejectOptions struct {
	KafkaEndpoints              []string
	KafkaTopic                  string
	KafkaConsumerGroupID        string
	TeamsWebhookNotificationUrl string
}

func toOptions() RejectOptions {
	return RejectOptions{
		KafkaEndpoints:              flags.GetKafkaEndpoints(),
		KafkaTopic:                  flags.GetKafkaTopic(),
		KafkaConsumerGroupID:        flags.GetKafkaConsumerGroupID(),
		TeamsWebhookNotificationUrl: flags.GetTeamsWebhookNotificationUrl(),
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject",
		Short: "Continuously report all rejected data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdutil.HandleError(toOptions().Run())
		},
	}

	addFlags(cmd)

	return cmd
}

func (o RejectOptions) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: o.KafkaEndpoints,
		Topic:   o.KafkaTopic,
		GroupID: o.KafkaConsumerGroupID,
	})

	for {
		message, err := dps.NextMessage(ctx, reader, dps.IsWebArchiveOwned)
		if err != nil {
			return fmt.Errorf("failed to read next message from kafka: %w", err)
		}

		slog.Info("Received reject message from DPS", "message", message.Response, "key", message.Key, "offset", message.Offset)

		if len(o.TeamsWebhookNotificationUrl) == 0 {
			continue
		}

		teamsMsg := teams.VerificationError(message, o.KafkaTopic, o.KafkaEndpoints)
		if err := teams.SendMessage(ctx, teamsMsg, o.TeamsWebhookNotificationUrl); err != nil {
			slog.Error("Failed to send message to teams", "error", err)
		}
	}
}
