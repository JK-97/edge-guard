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
	"io/ioutil"
	"jxcore/config"
	"jxcore/core"
	"jxcore/core/device"
	"jxcore/lowapi/ceph"
	"jxcore/lowapi/utils"
	"jxcore/subprocess"
	"jxcore/subprocess/gateway"
	"jxcore/version"
	"jxcore/web"
	"jxcore/web/route"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"

	// 调试
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	graceful = time.Second * 1
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

				// Network interface auto switch
				// Auto update /etc/resolv.conf to dnsmasq config
				// IoTEdge VPN auto reconnect
				errGroup.Go(func() error { return core.MaintainNetwork(ctx) })
			}

			if !noUpdate {
				log.Info("================Checking JxToolset Update===================")
				core.UpdateCore()
			}

			log.Info("================Configuring Environment===================")

			log.Info("ensure tmpfs is mounted")
			err = ceph.EnsureTmpFs()
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
