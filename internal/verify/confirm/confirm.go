package confirmImplmentation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hermetic/internal/dps"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func ReadConfirmTopic(ctx context.Context, reader *kafka.Reader, receiverUrl string) error {
	for {
		response, err := dps.ReadMessages(ctx, reader)
		if err != nil {
			return fmt.Errorf("failed to read message from confirm-topic: `%w`", err)
		}
		if len(receiverUrl) > 0 {
			err = SendConfirmMessage(receiverUrl, response.DPSResponse)
			if err != nil {
				return fmt.Errorf("failed to send confirm message: `%w`", err)
			}
		} else {
			fmt.Printf("Received message: %v\n", response.DPSResponse)
		}
	}
}

func SendConfirmMessage(baseUrl string, response dps.DigitalPreservationSystemResponse) error {
	slog.Info("Sending confirm message to receiver", "receiver_url", baseUrl)
	url, err := url.JoinPath(baseUrl, response.Identifier)
	if err != nil {
		return fmt.Errorf("failed to join URL path: `%w`", err)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal DPS response to JSON: `%w`", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: `%w`", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: `%w`", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	fmt.Printf("Successfully sent confirm message to %s\n", url)
	return nil
}
