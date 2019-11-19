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
	"context"
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"io/ioutil"
	"jxcore/config"
	"jxcore/core"
	"jxcore/core/device"
	"jxcore/lowapi/ceph"
	"jxcore/lowapi/dns"
	"jxcore/lowapi/utils"
	"jxcore/subprocess"
	"jxcore/subprocess/gateway"
	"jxcore/version"
	"jxcore/web"
	"jxcore/web/route"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	// 调试
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"golang.org/x/sync/errgroup"
)

const (
	graceful = time.Second * 15
)

var (
	debug    bool   = false
	port     string = ":80"
	noUpdate bool   = false
)

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
			log.Info("==================Jxcore Serve Starting=====================")

			log.Info("================Checking Edgenode Status===================")
			if !utils.Exists(InitPath) {
				log.Fatal("please run the bootstrap before serve")
			}
			currentdevice, err := device.GetDevice()
			if err != nil {
				log.Fatal(err)
			}
			log.WithFields(log.Fields{"INFO": "Device"}).Info("workerid : ", currentdevice.WorkerID)

			ctx, cancel := context.WithCancel(context.Background())
			errGroup, ctx := errgroup.WithContext(ctx)

			core := core.GetJxCore()
			if device.GetDeviceType() == version.Pro {
				log.Info("=======================Configuring Network============================")
				core.ConfigNetwork()
			}

			// Network interface auto switch
			// Auto update /etc/resolv.conf to dnsmasq config
			// IoTEdge VPN auto reconnect
			errGroup.Go(core.MaintainNetwork)

			if !noUpdate {
				log.Info("================Checking JxToolset Update===================")
				core.UpdateCore()
			}

			log.Info("================Configuring Environment===================")
			// ensure tmpfs is mounted
			err = ceph.EnsureTmpFs()
			if err != nil {
				log.Fatal(err)
			}
			// ensure docker start with correct dnsmasq setup
			err = ensureDocker()
			if err != nil {
				log.Fatal(err)
			}
			core.ConfigSupervisor()

			log.Info("================Starting Subprocesses===================")

			// start up all component process
			errGroup.Go(func() error { return subprocess.RunServer(ctx) })
			errGroup.Go(func() error { return subprocess.RunJxserving(ctx) })
			errGroup.Go(func() error { return subprocess.RunMcuProcess(ctx) })

			// web server
			errGroup.Go(gateway.ServeGateway)
			errGroup.Go(func() error { return web.Serve(ctx, port, route.Routes(), graceful) })
			errGroup.Go(func() error { return web.Serve(ctx, ":10880", http.DefaultServeMux, graceful) })

			// handle SIGTERM and SIGINT
			go func() {
				termChan := make(chan os.Signal, 1)
				signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
				sig := <-termChan
				log.WithFields(log.Fields{"signal": sig}).Info("receive a signal to stop all process & exit")
				cancel()

				time.Sleep(graceful)
				log.Info("===============Jxcore exited===============")
				log.Fatal("Jxcore cannot exit within graceful period: ", graceful.String())
			}()

			err = errGroup.Wait()
			log.Info("===============Jxcore exited===============")
			if err != nil {
				log.Error("Exited with error: %+v", err)
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&port, "port", port, "Port to run Application server on")
	serveCmd.PersistentFlags().BoolVar(&debug, "debug", debug, "Whether to enable pprof")
	serveCmd.PersistentFlags().BoolVar(&noUpdate, "no-update", noUpdate, "Whether to check for update")

	serveCmd.PersistentFlags().String("interface", "eth0", "gateway listen where")
	serveCmd.PersistentFlags().String("config", "./settings.yaml", "yaml setting for component")
	cfg := config.Config()
	_ = cfg.BindPFlag("yamlsettings", serveCmd.PersistentFlags().Lookup("config"))
	_ = cfg.BindPFlag("interface", serveCmd.PersistentFlags().Lookup("interface"))
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
func ensureDocker() error {
	exec.Command("sed", "-i", "s/.*172.17.0.1/#listen/", "/etc/dnsmasq.conf").Run()
	dns.RestartDnsmasq()

	if dockerNeedRestart() {
		if err := exec.Command("service", "docker", "restart").Run(); err != nil {
			return err
		}
	}

	if _, err := netlink.LinkByName("docker0"); err != nil {
		return err
	}
	if err := exec.Command("sed", "-i", "s/#listen/listen-address=172.17.0.1/", "/etc/dnsmasq.conf").Run(); err != nil {
		return err
	}
	dns.RestartDnsmasq()
	return nil
}

func dockerNeedRestart() bool {
	data, err := ioutil.ReadFile("/var/run/docker.pid")
	if err != nil {
		return true
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return true
	}
	_, err = os.FindProcess(pid)
	return err != nil
}
