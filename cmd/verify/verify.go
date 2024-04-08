package verify

import (
	"github.com/spf13/cobra"
	"hermetic/cmd/verify/confirm"
	"hermetic/cmd/verify/reject"
)

func NewCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "verify",
		Short: "Continuously verifies uploaded data responses",
	}
	rootCommand.AddCommand(reject.NewCommand())
	rootCommand.AddCommand(confirm.NewCommand())
	return rootCommand
}
