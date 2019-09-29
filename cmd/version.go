package cmd

import (
	"fmt"
	"jxcore/journal"
	"jxcore/version"

	"github.com/spf13/cobra"

	// 日志采插件
	_ "jxcore/journal/docker"
	_ "jxcore/journal/rfile"
	_ "jxcore/journal/systemd"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of jxcore",
	Long:  `All software has versions. This jxcore-backend`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		fmt.Println("Journal Modes:", journal.RegisteredModes())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
