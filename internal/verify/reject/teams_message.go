package rejectImplementation

import (
	"fmt"
	"hermetic/internal/dps"
	"hermetic/internal/teams"
	"strconv"
	"strings"
)

func createTeamsDigitalPreservationSystemFailureMessage(message *dps.KafkaResponse, rejectTopicName string, kafkaEndpoints []string) teams.Message {
	facts := []teams.Fact{
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
			Value: message.DPSResponse.Identifier,
		},
		{
			Name:  "Urn",
			Value: message.DPSResponse.Urn,
		},
		{
			Name:  "Path",
			Value: message.DPSResponse.Path,
		},
		{
			Name:  "ContentType",
			Value: message.DPSResponse.ContentType,
		},
		{
			Name:  "ContentCategory",
			Value: message.DPSResponse.ContentCategory,
		},
		{
			Name:  "Date of submission",
			Value: message.DPSResponse.Date,
		},
	}
	for index, check := range message.DPSResponse.Checks {
		facts = append(facts, teams.Fact{
			Name:  fmt.Sprintf("Check #%d status", index),
			Value: check.Status,
		})
		facts = append(facts, teams.Fact{
			Name:  fmt.Sprintf("Check #%d message", index),
			Value: check.Message,
		})
		facts = append(facts, teams.Fact{
			Name:  fmt.Sprintf("Check #%d reason", index),
			Value: check.Reason,
		})
		facts = append(facts, teams.Fact{
			Name:  fmt.Sprintf("Check #%d file", index),
			Value: check.File,
		})
	}

	return teams.Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    "Verification error",
		Sections: []teams.Section{
			{
				ActivityTitle:    "Verification error",
				ActivitySubtitle: "A Digital Preservation System (DPS) upload failed",
				ActivityImage:    "https://www.dictionary.com/e/wp-content/uploads/2018/03/thisisfine-1-300x300.jpg",
				Facts:            facts,
			},
		},
	}
}
