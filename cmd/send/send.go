package send

import (
	"fmt"
	"hermetic/internal/common_flags"
	sendImplementation "hermetic/internal/send"
	"hermetic/internal/teams"

	"github.com/spf13/cobra"
)

const (
	transferTopicFlagName      string = "transfer-topic"
	stageArtifactsRootFlagName string = "stage-artifacts-root"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "send",
		Short: "Continuously sends data to digital storage",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallSend,
	}
	command.Flags().String(transferTopicFlagName, "", "name of transfer-topic")
	if markTransferRequiredError := command.MarkFlagRequired(transferTopicFlagName); markTransferRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", transferTopicFlagName, markTransferRequiredError))
	}

	command.Flags().String(stageArtifactsRootFlagName, "", `path to the root directory with the following content structure:
/stage-artifacts-root
├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm
│   ├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
│   └── /checksum_transferred.md5
└── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm
    ├── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
    └── /checksum_transferred.md5
... etc`)
	if err := command.MarkFlagRequired(stageArtifactsRootFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", stageArtifactsRootFlagName, err))
	}

	return command
}

func parseArgumentsAndCallSend(cmd *cobra.Command, args []string) error {
	transferTopicName, err := cmd.Flags().GetString(transferTopicFlagName)
	if err != nil {
		return fmt.Errorf("getting transfer topic name failed, original error: '%w'", err)
	}
	stageArtifactsRoot, err := cmd.Flags().GetString("stage-artifacts-root")
	if err != nil {
		return fmt.Errorf("getting stage artifacts root failed, original error: '%w'", err)
	}

	err = sendImplementation.PrepareAndSendSubmissionInformationPackage(common_flags.KafkaEndpoints, transferTopicName, stageArtifactsRoot)
	if err != nil {
		err = fmt.Errorf("transfer error, cause: `%w`", err)
		fmt.Printf("Sending error message to Teams\n")
		teamsErrorMessage := teams.CreateGeneralFailureMessage(err)
		if err := teams.SendMessage(teamsErrorMessage, common_flags.TeamsWebhookNotificationUrl); err != nil {
			err = fmt.Errorf("failed to send error message to Teams, cause: `%w`", err)
			fmt.Printf("%s\n", err)
		}
		return err
	}
	return err
}
