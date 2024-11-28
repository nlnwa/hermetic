package teams

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nlnwa/hermetic/internal/dps"
)

func TestCreateTeamsMessage(t *testing.T) {
	kafkaResponse := &dps.KafkaMessage{
		Offset: 0,
		Key:    "key",
		Value: dps.Message{
			Date:            "date",
			Identifier:      "identifier",
			Urn:             "urn",
			Path:            "path",
			ContentType:     "contentType",
			ContentCategory: "contentCategory",
			Checks: []dps.Check{
				{
					Status:  "status",
					Message: "message",
					Reason:  "reason",
					File:    "file",
				},
			},
		},
	}

	rejectTopicName := "rejectTopicName"
	kafkaEndpoints := []string{"kafkaEndpoints"}

	message := VerificationError(
		kafkaResponse,
		rejectTopicName,
		kafkaEndpoints,
	)
	expectedMessage := Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    "Verification error",
		Sections: []Section{
			{
				ActivityTitle:    "Verification error",
				ActivitySubtitle: "A Digital Preservation System (DPS) upload failed",
				ActivityImage:    "https://www.dictionary.com/e/wp-content/uploads/2018/03/thisisfine-1-300x300.jpg",

				Facts: []Fact{
					{
						Name:  "Kafka message offset",
						Value: "0",
					},
					{
						Name:  "Kafka message key",
						Value: "key",
					},
					{
						Name:  "Kafka topic",
						Value: "rejectTopicName",
					},
					{
						Name:  "Kafka endpoints",
						Value: "kafkaEndpoints",
					},
					{
						Name:  "Identifier",
						Value: "identifier",
					},
					{
						Name:  "Urn",
						Value: "urn",
					},
					{
						Name:  "Path",
						Value: "path",
					},
					{
						Name:  "ContentType",
						Value: "contentType",
					},
					{
						Name:  "ContentCategory",
						Value: "contentCategory",
					},
					{
						Name:  "Date of submission",
						Value: "date",
					},
					{
						Name:  "Check #0 status",
						Value: "status",
					},
					{
						Name:  "Check #0 message",
						Value: "message",
					},
					{
						Name:  "Check #0 reason",
						Value: "reason",
					},
					{
						Name:  "Check #0 file",
						Value: "file",
					},
				},
			},
		},
	}
	if !cmp.Equal(message, expectedMessage) {
		t.Errorf("Expected message to be: '%v'\n, got: \n'%v'", prettify(expectedMessage), prettify(message))
	}

}

func prettify(message Message) string {
	s, err := json.MarshalIndent(message, "", "\t")

	if err != nil {
		panic(fmt.Sprintf("Failed to prettify message: %v", err))
	}
	return string(s)
}
