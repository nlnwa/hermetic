package acquisition

import (
	"fmt"
	acquisitionImplementation "hermetic/internal/acquisition"
	"hermetic/internal/common_flags"

	"github.com/spf13/cobra"
)

const (
	transferTopicFlagName   string = "transfer-topic"
	acquisitionRootFlagName string = "acquisition-root"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "acquisition",
		Short: "Uploads data to digital storage",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallAcquisition,
	}
	command.Flags().String(transferTopicFlagName, "", "name of transfer-topic")
	if markTransferRequiredError := command.MarkFlagRequired(transferTopicFlagName); markTransferRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", transferTopicFlagName, markTransferRequiredError))
	}

	command.Flags().String(acquisitionRootFlagName, "", `path to the root directory with the following content structure:
/acquisition-root
├── checksums.md5
├── checksum_transferred.md5
├── /acquisition.yaml
├── <other-small-and-few-files>
└── /<other-files-and-directories>
`)
	if err := command.MarkFlagRequired(acquisitionRootFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", acquisitionRootFlagName, err))
	}

	return command
}

func parseArgumentsAndCallAcquisition(cmd *cobra.Command, args []string) error {
	transferTopicName, err := cmd.Flags().GetString(transferTopicFlagName)
	if err != nil {
		return fmt.Errorf("getting transfer topic name failed, original error: '%w'", err)
	}
	acquisitionRoot, err := cmd.Flags().GetString(acquisitionRootFlagName)
	if err != nil {
		return fmt.Errorf("getting stage artifacts root failed, original error: '%w'", err)
	}

	err = acquisitionImplementation.PrepareAndSendSubmissionInformationPackage(common_flags.KafkaEndpoints, transferTopicName, acquisitionRoot)
	if err != nil {
		return fmt.Errorf("transfer error, cause: `%w`", err)
	}

	return nil
}
