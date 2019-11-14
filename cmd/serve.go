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
	"io/ioutil"
	"jxcore/config"
	"jxcore/core"
	"jxcore/core/device"
	log "jxcore/go-utils/logger"
	"jxcore/lowapi/ceph"
	"jxcore/lowapi/dns"
	"jxcore/lowapi/utils"
	"jxcore/subprocess"
	"jxcore/subprocess/gateway"
	"jxcore/version"
	"jxcore/web/route"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"

	// 调试
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

var start chan bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve http backend for jxcore",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Deamonize(func() {
			c := exec.Command("sed", "-i", "s/.*172.17.0.1/#listen/", "/etc/dnsmasq.conf")
			c.Run()
			dns.RestartDnsmasq()
			ceph.CheckTmpFs()
			core := core.GetJxCore()
			go func() {
				gateway.Setup()
				gateway.ServeGateway()
			}()
			forever := make(chan interface{}, 1)

			if utils.Exists(InitPath) {
			} else {
				log.Fatal("please run the bootstrap before serve")
			}
			currentdevice, err := device.GetDevice()
			utils.CheckErr(err)
			log.WithFields(log.Fields{"INFO": "Device"}).Info("workerid : ", currentdevice.WorkerID)

			go subprocess.RunMcuProcess()

			once := &sync.Once{}
			if device.GetDeviceType() == version.Pro {
				core.ConfigNetwork()
			}

			ensureDocker()
			go once.Do(func() {
				subprocess.RunJxserving()
			})

			flags := cmd.Flags()

			if noUpdate, _ := flags.GetBool("no-update"); !noUpdate {
				core.UpdateCore()
			}
			core.ConfigSupervisor()

			//collection log
			if _, err = os.Stat(LogsPath); err != nil {
				exec.Command("mkdir", "-p", LogsPath)
			}
			//core.CollectJournal(currentdevice.WorkerID)

			//start up all component process
			go subprocess.Run()
			log.Info("all process has run")

			//web server
			port, err := flags.GetString("port")
			if err != nil {
				port = ":80"
			}
			go func() {
				log.Info("Listen on", port)
				log.Fatal(http.ListenAndServe(port, route.Routes()))
				os.Exit(1)
				forever <- nil
			}()
			if debug, _ := flags.GetBool("debug"); debug {
				go func() {
					port := ":10880"
					log.Info("Enable Debug Mode Listen on", port)
					log.Fatal(http.ListenAndServe(port, nil))
					os.Exit(1)
					forever <- nil
				}()
			}

			<-forever
		})

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().String("port", ":80", "Port to run Application server on")
	serveCmd.PersistentFlags().String("interface", "eth0", "gateway listen where")
	serveCmd.PersistentFlags().String("config", "./settings.yaml", "yaml setting for component")
	serveCmd.PersistentFlags().Bool("debug", false, "Whether to enable pprof")
	serveCmd.PersistentFlags().Bool("no-update", false, "Whether to check for update")
	cfg := config.Config()
	cfg.BindPFlag("yamlsettings", serveCmd.PersistentFlags().Lookup("config"))
	cfg.BindPFlag("interface", serveCmd.PersistentFlags().Lookup("interface"))

}

// applySyncTools 配置同步工具
func applySyncTools() {
	if utils.Exists("/edge/synctools.zip") {
		data, err := ioutil.ReadFile("/edge/synctools.zip")
		if err != nil {
			log.Error(err)
		} else {
			err = utils.Unzip(data, "/edge/mnt")
			if err == nil {
				log.Info("has find the synctools.zip")
				os.Remove("/edge/synctools.zip.old")
				if err = os.Rename("/edge/synctools.zip", "/edge/synctools.zip.old"); err != nil {
					log.Error("Fail to move /edge/synctools.zip to /edge/synctools.zip.old", err)
				}
			}
		}
	}
}

// ensureDocker 确保 docker 服务会启动
func ensureDocker() {
	var err error
	data, err := ioutil.ReadFile("/var/run/docker.pid")
	if err == nil {
		pid, err := strconv.Atoi(string(data))
		if err == nil {
			_, err = os.FindProcess(pid)
		}
	}
	if err != nil {
		cmd := exec.Command("service", "docker", "restart")
		cmd.Run()
	}

	if _, err = netlink.LinkByName("docker0"); err == nil {
		cmd := exec.Command("sed", "-i", "s/#listen/listen-address=172.17.0.1/", "/etc/dnsmasq.conf")
		cmd.Run()
		dns.RestartDnsmasq()
	}
}
