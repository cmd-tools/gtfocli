package cmd

import (
	"github.com/cmd-tools/gtfocli/constants"
	"github.com/cmd-tools/gtfocli/logger"
	"github.com/spf13/cobra"
	"os"
)

var isDebugEnabled bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.Main,
	Short: "GFTO Command Line Interface",
	Long: `GFTO Command Line Interface for easy binaries search commands
that can be used to bypass local security restrictions in misconfigured systems.`,
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(constants.Main, isDebugEnabled)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&isDebugEnabled, "debug", false, "Enable application logging for debugging purposes.")
}

func IsDebug() bool {
	debug, _ := rootCmd.Flags().GetBool("debug")
	return debug
}
