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
	"jxcore/config"
	"jxcore/config/yaml"
	"jxcore/core"
	"jxcore/core/device"
	"jxcore/gateway"
	"jxcore/internal/network/ssdp"
	"jxcore/lowapi/ceph"
	"jxcore/subprocess"
	"jxcore/version"
	"jxcore/web"
	"jxcore/web/route"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "jxcore/lowapi/logger"

	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	graceful = time.Second * 15
)

var (
	debug    bool   = false
	addr     string = ":80"
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
			log.Infof("Config: %+v", yaml.Config)

			currentdevice, err := device.GetDevice()
			if err != nil {
				log.Fatal(err)
			}
			log.Info("workerid : ", currentdevice.WorkerID)

			ctx, cancel := context.WithCancel(context.Background())
			errGroup, ctx := errgroup.WithContext(ctx)

			ssdpClient := ssdp.NewClient(currentdevice.WorkerID, 5)
			errGroup.Go(func() error { return ssdpClient.Aliving(ctx) })
			if device.GetDeviceType() == version.Pro {
				log.Info("=======================Configuring Network============================")
				core.ConfigNetwork()

				// Network interface auto switch
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

			log.Info("================Starting Subprocesses===================")
			core.ConfigSupervisor()

			// start up all component process
			errGroup.Go(func() error { return subprocess.RunServer(ctx) })

			// web server
			errGroup.Go(func() error { return gateway.ServeGateway(ctx, graceful) })
			errGroup.Go(func() error { return web.Serve(ctx, addr, route.Routes(), graceful) })
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
	serveCmd.PersistentFlags().StringVar(&addr, "port", addr, "Addr to run Application server on")
	serveCmd.PersistentFlags().BoolVar(&debug, "debug", debug, "Whether to enable pprof")
	serveCmd.PersistentFlags().BoolVar(&noUpdate, "no-update", noUpdate, "Whether to check for update")

	serveCmd.PersistentFlags().String("interface", "eth0", "gateway listen where")
	serveCmd.PersistentFlags().String("config", "./settings.yaml", "yaml setting for component")
	cfg := config.Config()
	_ = cfg.BindPFlag("yamlsettings", serveCmd.PersistentFlags().Lookup("config"))
	_ = cfg.BindPFlag("interface", serveCmd.PersistentFlags().Lookup("interface"))
}
