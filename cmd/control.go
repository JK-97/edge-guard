package cmd

import (
	"github.com/JK-97/edge-guard/ctl/client"
	"github.com/JK-97/edge-guard/version"

	"github.com/spf13/cobra"
)

// controlCmd represents the control command
var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.Run(version.Version)
	},
}

func init() {
	rootCmd.AddCommand(controlCmd)
}
