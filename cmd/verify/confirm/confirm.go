package confirm

import (
	"fmt"
	"github.com/spf13/cobra"
	"hermetic/internal/common_flags"
	"hermetic/internal/teams"
	confirmImplementation "hermetic/internal/verify/confirm"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "confirm",
		Short: "Continuously report all successfully preserved data",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallVerify,
	}
	confirmTopicFlagName := "confirm-topic"
	command.Flags().String(confirmTopicFlagName, "", "name of confirm-topic")
	if err := command.MarkFlagRequired(confirmTopicFlagName); err != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", confirmTopicFlagName, err))
	}
	return command
}

func parseArgumentsAndCallVerify(cmd *cobra.Command, args []string) error {
	confirmTopicName, err := cmd.Flags().GetString("confirm-topic")
	if err != nil {
		return fmt.Errorf("failed to get confirm-topic flag, cause: `%w`", err)
	}

	err = confirmImplementation.Verify(confirmTopicName)
	if err != nil {
		err = fmt.Errorf("verification error, cause: `%w`", err)
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
