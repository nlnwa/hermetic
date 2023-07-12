package verify

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "verify",
		Short: "Continuously verifies uploaded data responses",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallVerify,
	}
	rejectTopicFlagName := "reject-topic"
	command.Flags().String(rejectTopicFlagName, "", "name of reject-topic")
	if err := command.MarkFlagRequired(rejectTopicFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", rejectTopicFlagName, err))
	}
	// TODO(https://github.com/nlnwa/hermetic/issues/3): handle the `confirm`
	// topic messages
	return command
}

func parseArgumentsAndCallVerify(cmd *cobra.Command, args []string) error {
	fmt.Println("Parsing verify arguments")
	return verify()
}

func verify() error {
	fmt.Println("Running verify")
	return nil
}


package verifier

import (
	"_warc-to-storage/internal/endpoints"
	"_warc-to-storage/internal/flag"
	kafkaTopics "_warc-to-storage/internal/topics"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	confirmTopicName string
	rejectTopicName  string
	kafkaEndpoints   []string
}

func NewCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "verify",                                        // TODO
		Short: "Continuously verifies uploaded data responses", // TODO
		Long:  ``,                                              // TODO
		Args:  cobra.NoArgs,
		RunE:  entryPointCommand,
	}

	cmd.Flags().String(flag.ConfirmTopicName, "", "confirm topic") // TODO
	cmd.MarkFlagRequired(flag.ConfirmTopicName)
	cmd.Flags().String(flag.RejectTopicName, "", "reject topic") // TODO
	cmd.MarkFlagRequired(flag.RejectTopicName)
	flag.AddCommonFlags(cmd)
	//cmd.Flags().StringSlice(flag.KafkaEndpoints, []string{}, "list of kafka endpoints")
	//cmd.MarkFlagRequired(flag.KafkaEndpoints)
	//cmd.Flags().StringP(flag.TargetEnvironment, "t", "_", "what environment to target, most commonly `stage` or `prod`")
	//cmd.MarkFlagRequired(flag.TargetEnvironment)

	return cmd
}

func entryPointCommand(cmd *cobra.Command, args []string) error {
	c := &config{}
	c.confirmTopicName = viper.GetString(flag.ConfirmTopicName)
	c.rejectTopicName = viper.GetString(flag.RejectTopicName)
	c.kafkaEndpoints = viper.GetStringSlice(flag.KafkaEndpoints)
	fmt.Println(viper.GetString("rubbish"))
	fmt.Println(args)
	fmt.Println(cmd)

	return verifyDataSentToDigitalStorage(c)
}

func verifyDataSentToDigitalStorage(myConfig *config) error { //(rejectTopicName string, confirmTopicName string, kafkaEndpoints []string) error {
	fmt.Printf("This is reject topic: %s\n", myConfig.rejectTopicName)
	fmt.Printf("This is confirm topic: %s\n", myConfig.confirmTopicName)
	fmt.Printf("This is kafka endpoints: %s\n", myConfig.kafkaEndpoints)
	return nil
}

func verifyDataSentToDigitalStorageDummy() error {
	hosts := endpoints.Stage

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: hosts,
		Topic:   kafkaTopics.Reject,
		//Partition: 0,return nil
		//MaxBytes: 10e6, // 10MB
	})
	//r.SetOffset(42)

	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			//break
			return fmt.Errorf("something bad happened, cause: `%w`", err)
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
	return nil

	//if err := r.Close(); err != nil {
	//	log.Fatal("failed to close reader:", err)
	//}

}

//func init() {
//	rootCmd.AddCommand(reverseCmd)
//}
