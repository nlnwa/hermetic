package send

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "send",
		Short: "Continuously sends data to digital storage",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallSend,
	}
	transferTopicFlagName := "transfer-topic"
	command.Flags().String(transferTopicFlagName, "", "name of transfer-topic")
	if markTransferRequiredError := command.MarkFlagRequired(transferTopicFlagName); markTransferRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", transferTopicFlagName, markTransferRequiredError))
	}

	rootPathFlagName := "stage-artifacts-root"
	command.Flags().String(rootPathFlagName, "", `path to the root directory with the following content structure:
/stage-artifacts-root
├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm
│   ├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
│   └── /checksum_transferred.md5
└── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm
    ├── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
    └── /checksum_transferred.md5
... etc`)
	if markRootPathRequiredError := command.MarkFlagRequired(rootPathFlagName); markRootPathRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", rootPathFlagName, markRootPathRequiredError))
	}

	return command
}

func parseArgumentsAndCallSend(cmd *cobra.Command, args []string) error {
	fmt.Println("Parsing send arguments")
	return send()
}

func send() error {
	fmt.Println("Running send")
	return nil
}
