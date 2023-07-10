package common

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ParseKafkaEndpointsFlags(cmd *cobra.Command) ([]string, error) {
	kafkaEndpoints, kafkaEndpointsError := cmd.Flags().GetStringSlice("kafka-endpoints")
	if kafkaEndpointsError != nil {
		return []string{}, fmt.Errorf("failed to get kafka-endpoints flag, cause: `%w`", kafkaEndpointsError)
	}
	return kafkaEndpoints, nil
}
