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
	"bufio"
	"fmt"
	"jxcore/app/model/docker"
	"jxcore/log"
	"jxcore/regeister"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	vpnmode string

	ticket string

	authHost string

	skipRestore bool
)

const (
	restoreImagePath     = "/edge/jxbootstrap/worker/dependencies/recover/dockerimage"
	restoreBootstrapPath = "/edge/jxbootstrap"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap http backend for jxcore",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		workerid := regeister.BuildWokerID()

		if ticket == "" {
			fmt.Println("Need Ticket")
			fmt.Println("Worker ID:", workerid)
			fmt.Println("Please enter ticket:")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			ticket = scanner.Text()
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			}
		}
		if len(ticket) < 2 {
			fmt.Fprintln(os.Stderr, "Wrong Ticket. Too short:", ticket)
			return
		}

		if !skipRestore {
			if _, err := os.Stat(restoreImagePath); err == nil {
				os.RemoveAll("/restore/dockerimage")
				os.Mkdir("/restore", 0644)
				err = os.Rename(restoreImagePath, "/restore/dockerimage")
				if err != nil {
					log.Error(err)
				}
			}

			log.Info("Restore Docker Images")
			var dockerobj = docker.NewClient()
			err := dockerobj.DockerRestore()
			if err != nil {
				log.Error(err)
			} else {
				log.Info("Finish Restore Docker Images")
			}

			err = exec.Command("hostnamectl", "set-hostname", "worker-"+workerid).Run()
			if err != nil {
				panic(err)
			}

			if _, err := os.Stat(restoreBootstrapPath); err == nil {
				os.RemoveAll("/jxbootstrap")
				err = os.Rename(restoreBootstrapPath, "/jxbootstrap")
				if err != nil {
					log.Error(err)
				}
			}

			basecmd := exec.Command("/jxbootstrap/worker/scripts/base.sh")
			basecmd.Stdout = os.Stdout
			basecmd.Stdout = os.Stderr
			err = basecmd.Run()
			if err != nil {
				panic(err)
			}

			//LookUpDns
			time.Sleep(3 * time.Second)

		}
		if authHost == "" {
			authHost = regeister.FallBackAuthHost
		}

		host := GetHost(authHost)

		// "auth.iotedge.jiangxingai.com"
		regeister.LookUpDns(host)

		initcmd := exec.Command("touch", "/edge/init")
		initcmd.Run()

		log.Info("Register to ", authHost)
		regeister.BuildDeviceInfo(vpnmode, ticket, authHost)

		// exec.Command("rm", "-r", "/jxbootstrap").Start()
	},
}

// GetHost 从 url 中解析 Host
func GetHost(u string) string {
	uri, err := url.Parse(u)
	if err != nil {
		return u
	}
	return uri.Hostname()
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.AddCommand(bootstrapCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	bootstrapCmd.PersistentFlags().StringVarP(&vpnmode, "mode", "m", regeister.VPNModeRandom, "openvpn or wireguard or local")
	bootstrapCmd.PersistentFlags().StringVarP(&ticket, "ticket", "t", "", "ticket for bootstrap")
	bootstrapCmd.PersistentFlags().StringVarP(&authHost, "host", "", regeister.FallBackAuthHost, "host for bootstrap")
	bootstrapCmd.PersistentFlags().BoolVarP(&skipRestore, "skip", "s", false, "skip restore")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
