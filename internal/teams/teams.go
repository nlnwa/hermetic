package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nlnwa/hermetic/internal/dps"
)

const (
	avoidMicrosoftTeamsWebhookRateLimit = 1 * time.Second
)

type Fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Section struct {
	ActivityTitle    string `json:"activityTitle"`
	ActivitySubtitle string `json:"activitySubtitle"`
	ActivityImage    string `json:"activityImage"`
	Facts            []Fact `json:"facts"`
}

type Message struct {
	Type       string    `json:"@type"`
	Context    string    `json:"@context"`
	ThemeColor string    `json:"themeColor"`
	Summary    string    `json:"summary"`
	Sections   []Section `json:"sections"`
}

func SendMessage(ctx context.Context, payload Message, webhookUrl string) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal teams message: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookUrl, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to teams: %w", err)
	}
	defer resp.Body.Close()

	time.Sleep(avoidMicrosoftTeamsWebhookRateLimit)

	return nil
}

func Error(err error) Message {
	return Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    "System error",
		Sections: []Section{
			{
				ActivityTitle:    "System error",
				ActivitySubtitle: "A Digital Preservation System (DPS) general failure",
				ActivityImage:    "https://www.dictionary.com/e/wp-content/uploads/2018/03/thisisfine-1-300x300.jpg",
				Facts: []Fact{
					{
						Name:  "Error",
						Value: err.Error(),
					},
				},
			},
		},
	}
}

func VerificationError(message *dps.KafkaMessage, rejectTopicName string, kafkaEndpoints []string) Message {
	facts := []Fact{
		{
			Name:  "Kafka message offset",
			Value: strconv.FormatInt(message.Offset, 10),
		},
		{
			Name:  "Kafka message key",
			Value: message.Key,
		},
		{
			Name:  "Kafka topic",
			Value: rejectTopicName,
		},
		{
			Name:  "Kafka endpoints",
			Value: strings.Join(kafkaEndpoints, ", "),
		},
		{
			Name:  "Identifier",
			Value: message.Value.Identifier,
		},
		{
			Name:  "Urn",
			Value: message.Value.Urn,
		},
		{
			Name:  "Path",
			Value: message.Value.Path,
		},
		{
			Name:  "ContentType",
			Value: message.Value.ContentType,
		},
		{
			Name:  "ContentCategory",
			Value: message.Value.ContentCategory,
		},
		{
			Name:  "Date of submission",
			Value: message.Value.Date,
		},
	}
	for index, check := range message.Value.Checks {
		facts = append(facts, Fact{
			Name:  fmt.Sprintf("Check #%d status", index),
			Value: check.Status,
		})
		facts = append(facts, Fact{
			Name:  fmt.Sprintf("Check #%d message", index),
			Value: check.Message,
		})
		facts = append(facts, Fact{
			Name:  fmt.Sprintf("Check #%d reason", index),
			Value: check.Reason,
		})
		facts = append(facts, Fact{
			Name:  fmt.Sprintf("Check #%d file", index),
			Value: check.File,
		})
	}

	return Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    "Verification error",
		Sections: []Section{
			{
				ActivityTitle:    "Verification error",
				ActivitySubtitle: "A Digital Preservation System (DPS) upload failed",
				ActivityImage:    "https://www.dictionary.com/e/wp-content/uploads/2018/03/thisisfine-1-300x300.jpg",
				Facts:            facts,
			},
		},
	}
}
