package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	kafkaConsumerGroupIDFlagName string = "kafka-consumer-group-id"
	kafkaConsumerGroupIDHelp     string = "consumer group ID for Kafka consumer"
)

func AddKafkaFlags(cmd *cobra.Command) {
	cmd.Flags().String(kafkaConsumerGroupIDFlagName, "", kafkaConsumerGroupIDHelp)
}

func GetKafkaConsumerGroupID() string {
	return viper.GetString(kafkaConsumerGroupIDFlagName)
}
