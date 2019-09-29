package cmd

import (
	"jxcore/app/model/docker"
	"jxcore/app/model/mongo"
	"jxcore/app/model/pythonpackage"
	"jxcore/log"

	"github.com/spf13/cobra"
)

// RecoverCmd represents the bootstrap command
var RecoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "clean http backend for jxcore",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("docker"); v {
			log.Info("Restore Docker Images")
			var dockerobj = docker.NewClient()
			err := dockerobj.DockerRestore()
			if err != nil {
				log.Error(err)
			} else {
				log.Info("Finish Restore Docker Images")
			}
		}
		if v, _ := cmd.Flags().GetBool("python"); v {
			log.Info("Restore Python packages")
			var c = pythonpackage.NewPkgClient()
			c.RestorePyPkg()
			log.Info("Finish Restore Python packages")
		}
		if v, _ := cmd.Flags().GetBool("mongo"); v {
			log.Info("Restore MongoDB")
			mongo.InstallMongo()
			log.Info("Finish Restore MongoDB")
		}
	}}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.AddCommand(RecoverCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RecoverCmd.Flags().BoolP("docker", "d", true, "Recover docker image")
	RecoverCmd.Flags().BoolP("python", "p", true, "Recover python")
	RecoverCmd.Flags().BoolP("mongo", "m", false, "Recover Mongo")
}
