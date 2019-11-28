// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"jxcore/lowapi/docker"
	"jxcore/lowapi/pythonpkg"
	"os/exec"

	"github.com/spf13/cobra"
)

var ifall = "true"
var ifdocker = "false"
var ifpython = "false"

// bootstrapCmd represents the bootstrap command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean http backend for jxcore",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ifall == "true" {
			ifdocker = "true"
			ifpython = "true"
		}
		if ifdocker == "true" {
			var dockerobj = docker.NewClient()
			dockerobj.ContainerAllRemove()
			dockerobj.ImageAllRemove()
		}
		if ifpython == "true" {
			var c = pythonpkg.NewPkgClient()
			err := c.DeletePyPkg()
			if err != nil {
				log.Error(err)
			}
		}
		exec.Command("/bin/bash", "-c", "rm -r /tmp").Output()
		exec.Command("/bin/bash", "-c", "rm /edge/init").Output()

		exec.Command("fusermount", "-uz", "/data/edgebox/remote/").Output()
		exec.Command("rm", "/edge/init").Start()
	}}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.AddCommand(cleanCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	cleanCmd.PersistentFlags().StringVarP(&ifall, "all", "a", "true", "openvpn or wireguard or local")
	cleanCmd.PersistentFlags().StringVarP(&ifdocker, "docker", "d", "false", "openvpn or wireguard or local")
	cleanCmd.PersistentFlags().StringVarP(&ifpython, "python", "p", "false", "ticket for bootstrap")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
