package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlnwa/hermetic/cmd/acquisition"
	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/cmd/send"
	"github.com/nlnwa/hermetic/cmd/verify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hermetic",
		Short:         "hermetic - sends and verifies data for digital storage",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetEnvPrefix("HERMETIC")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			if err := loadConfig(); err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddGlobalFlags(cmd)

	cmd.AddCommand(send.NewCommand())
	cmd.AddCommand(verify.NewCommand())
	cmd.AddCommand(acquisition.NewCommand())
	return cmd
}

func loadConfig() error {
	if viper.IsSet("config") {
		// Read config file specified by 'config' flag
		viper.SetConfigFile(viper.GetString("config"))
		return viper.ReadInConfig()
	}

	// Read config file from default locations
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	// current directory
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	viper.AddConfigPath(workingDirectory)

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}

	return nil
}
