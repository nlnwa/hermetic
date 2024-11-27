package confirm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os/signal"
	"syscall"
	"time"

	"github.com/nlnwa/hermetic/cmd/internal/cmdutil"
	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/internal/dps"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	receiverUrlFlagName        string = "confirm-message-receiver"
	receiverUrlFlagHelpMessage string = "optional URL for confirm message receiver"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().String(receiverUrlFlagName, "", receiverUrlFlagHelpMessage)

	flags.AddKafkaFlags(cmd)
}

func toOptions() *ConfirmOptions {
	return &ConfirmOptions{
		KafkaTopic:           flags.GetKafkaTopic(),
		KafkaEndpoints:       flags.GetKafkaEndpoints(),
		KafkaConsumerGroupID: flags.GetKafkaConsumerGroupID(),
		ReceiverUrl:          viper.GetString(receiverUrlFlagName),
	}
}

type ConfirmOptions struct {
	KafkaTopic           string
	KafkaEndpoints       []string
	KafkaConsumerGroupID string
	ReceiverUrl          string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm",
		Short: "Continuously report all successfully preserved data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdutil.HandleError(toOptions().Run())
		},
	}

	addFlags(cmd)

	return cmd
}

func (o *ConfirmOptions) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
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

		slog.Info("Received confirm message from DPS", "message", message.Value, "key", message.Key, "offset", message.Offset)

		if len(o.ReceiverUrl) == 0 {
			continue
		}

		err = sendConfirmMessage(ctx, o.ReceiverUrl, message.Value)
		if err != nil {
			return fmt.Errorf("failed to send confirm message: %w", err)
		}
	}
}

func sendConfirmMessage(ctx context.Context, baseUrl string, response dps.Message) error {
	url, err := url.JoinPath(baseUrl, response.Identifier)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal DPS response to JSON: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: `%w`", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: `%w`", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
