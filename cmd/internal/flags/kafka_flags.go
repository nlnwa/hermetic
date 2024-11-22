package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	consumerGroupIDFlagName string = "consumer-group-id"
	consumerGroupIDHelp     string = "consumer group ID for Kafka consumer"
	defaultConsumerGroupID  string = "nettarkivet-hermetic-verify-confirm"
)

func AddKafkaFlags(cmd *cobra.Command) {
	cmd.Flags().String(consumerGroupIDFlagName, defaultConsumerGroupID, consumerGroupIDHelp)
}

func GetKafkaConsumerGroupID() string {
	return viper.GetString(consumerGroupIDFlagName)
}
