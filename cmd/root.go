package cmd

import (
	"hermetic/cmd/send"
	"hermetic/cmd/verify"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "hermetic",
		Short: "hermetic - sends and verifies data for digital storage",
	}
	rootCommand.PersistentFlags().StringSlice("kafka-endpoints", []string{}, "list of kafka endpoints")
	rootCommand.AddCommand(send.NewCommand())
	rootCommand.AddCommand(verify.NewCommand())
	return rootCommand
}
