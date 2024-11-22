package flags

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	kafkaEndpointsFlagName              string = "kafka-endpoints"
	kafkaTopicFlagName                  string = "kafka-topic"
	teamsWebhookNotificationUrlFlagName string = "teams-webhook-notification-url"
)

func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(kafkaTopicFlagName, "", "name of kafka topic")
	if err := cmd.MarkPersistentFlagRequired(kafkaTopicFlagName); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().StringSlice(kafkaEndpointsFlagName, []string{}, "list of kafka endpoints")
	if err := cmd.MarkPersistentFlagRequired(kafkaEndpointsFlagName); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().String(teamsWebhookNotificationUrlFlagName, "", "url to teams webhook for notifications")
}

func ValidateGlobalFlags() (err error) {
	if GetKafkaTopic() == "" {
		err = errors.New("kafka topic is required")
	}
	if len(GetKafkaEndpoints()) == 0 {
		err = errors.New("kafka endpoints are required")
	}
	return
}

func GetKafkaTopic() string {
	return viper.GetString(kafkaTopicFlagName)
}

func GetKafkaEndpoints() []string {
	return viper.GetStringSlice(kafkaEndpointsFlagName)
}

func GetTeamsWebhookNotificationUrl() string {
	return viper.GetString(teamsWebhookNotificationUrlFlagName)
}
