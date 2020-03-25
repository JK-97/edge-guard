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
	"fmt"
	"jxcore/internal/network/dns"
	"jxcore/lowapi/docker"
	"jxcore/lowapi/utils"
	"net/http"
	"strings"

	// 调试
	_ "net/http/pprof"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func SetColor(msg string, conf, bg, text int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, text, msg, 0x1B)
}

// serveCmd represents the serve command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check jxcore runtime",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(SetColor("Step One\n", 0, 0, TextBlue))
		fmt.Print(SetColor("**********Check dependent software***********\n", 0, 0, TextBlue))
		for _, Denpends := range DependsFile {
			if _, err := exec.LookPath(Denpends); err == nil {
				fmt.Print(SetColor("INSTALLED   "+Denpends+"\n", 0, 0, TextGreen))
			} else {
				fmt.Print(SetColor("UNINSTALLED "+Denpends+"\n", 0, 0, TextYellow))
			}

		}
		fmt.Print(SetColor("Step Two\n", 0, 0, TextBlue))
		fmt.Print(SetColor("****************Check the SDK****************\n", 0, 0, TextBlue))
		out, err := exec.Command("/bin/bash", "-c", "export LC_ALL=C && pip3 freeze | grep edgebox").Output()
		utils.CheckErr(err)
		if strings.Contains(string(out), "edgebox==") {
			fmt.Print(SetColor("INSTALLED   edgenode sdk"+"\n", 0, 0, TextGreen))
		} else {
			fmt.Print(SetColor("UNINSTALLED edgenode sdk"+"\n", 0, 0, TextYellow))
		}
		fmt.Print(SetColor("Step Three\n", 0, 0, TextBlue))
		fmt.Print(SetColor("*********Check the docker container**********\n", 0, 0, TextBlue))
		dockerImages, err := docker.ImagesList()
		dockerImagesTag := make([]string, 0)
		for _, dockerimage := range dockerImages {
			dockerImagesTag = append(dockerImagesTag, dockerimage.Tag[0])
		}

		for _, dependimage := range DependsImages {
			flag := false
			for _, imagetag := range dockerImagesTag {
				if dependimage == imagetag {
					fmt.Print(SetColor("INSTALLED   "+dependimage+"\n", 0, 0, TextGreen))
					flag = true
					break
				}
			}
			if flag == false {
				fmt.Print(SetColor("UNINSTALLED "+dependimage+"\n", 0, 0, TextYellow))
			}

		}

		fmt.Print(SetColor("Step Four\n", 0, 0, TextBlue))
		fmt.Print(SetColor("************Check network status*************\n", 0, 0, TextBlue))
		err = exec.Command("ping", "baidu.com", "-c", "1", "-W", "5").Run()
		if err != nil {
			fmt.Print(SetColor("NETWORKER CHECK STATUS BAD \n", 0, 0, TextRed))
		} else {
			fmt.Print(SetColor("NETWORKER CHECK STATUS OK \n", 0, 0, TextGreen))
		}

		fmt.Print(SetColor("Step Five\n", 0, 0, TextBlue))
		fmt.Print(SetColor("***************Check Dns config**************\n", 0, 0, TextBlue))
		if err := dns.CheckDnsmasqConf(); err == nil {
			fmt.Print(SetColor("DNS CONFIG CHECK OK \n", 0, 0, TextGreen))
		} else {
			fmt.Print(SetColor("DNS CONFIG CHECK BAD \n", 0, 0, TextRed))
			fmt.Print(err)
		}
		if out, _ := exec.Command("pgrep", "dnsmasq").Output(); out != nil {
			fmt.Print(SetColor("RUNNING     "+"dnsmasq"+"\n", 0, 0, TextGreen))
		} else {
			fmt.Print(SetColor("dnsmasq RUNNING"+"\n", 0, 0, TextYellow))
		}

		fmt.Print(SetColor("Step Six\n", 0, 0, TextBlue))
		fmt.Print(SetColor("******Check component operation status********\n", 0, 0, TextBlue))
		if out, _ := exec.Command("/bin/bash", "-c", "ps -ef | grep \"jxcore serve\" | grep -v grep | awk '{print $2}'").Output(); len(out) != 0 {
			fmt.Print(SetColor("has detect the jxcore running"+"\n", 0, 0, TextGreen))
			resp, err := http.Get("http://localhost:80/ping")
			utils.CheckErr(err)
			if resp.StatusCode == 200 {
				fmt.Print(SetColor("jxcore web service health"+"\n", 0, 0, TextGreen))
			}
			for programname, shellpath := range BinFilesMAP {
				if out, _ := exec.Command("pgrep", "-f", shellpath).Output(); len(out) != 0 {
					fmt.Print(SetColor("RUNNING     "+programname+"\n", 0, 0, TextGreen))
				} else {

					fmt.Print(SetColor("NOT RUNNING "+programname+"\n", 0, 0, TextYellow))
				}
			}
		} else {
			fmt.Print(SetColor("JXCORE NOT RUNNING"+"\n", 0, 0, TextYellow))
		}

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

}
