package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carlmjohnson/requests"
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

func SendMessage(payload Message, webhookUrl string) error {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Sending message to Teams")
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message, cause: `%w`", err)
	}
	err = requests.
		URL(webhookUrl).
		BodyBytes(bytes).
		ContentType("text/plain").
		Fetch(timeoutContext)
	if err != nil {
		return fmt.Errorf("failed to send message to teams, cause: `%w`", err)
	}
	return nil
}

func CreateGeneralFailureMessage(failure error) Message {
	facts := []Fact{
		{
			Name:  "Error",
			Value: failure.Error(),
		},
	}
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
				Facts:            facts,
			},
		},
	}
}
