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
