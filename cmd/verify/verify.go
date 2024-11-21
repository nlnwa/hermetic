package verify

import (
	"github.com/nlnwa/hermetic/cmd/verify/confirm"
	"github.com/nlnwa/hermetic/cmd/verify/reject"
	"github.com/spf13/cobra"
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
