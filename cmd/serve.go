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
	"os"
	"os/signal"
	"syscall"
	"time"

	"jxcore/config"
	"jxcore/core"
	"jxcore/core/device"
	"jxcore/gateway"
	"jxcore/internal/network/ssdp"
	"jxcore/lowapi/ceph"
	"jxcore/lowapi/logger"
	"jxcore/monitor"
	"jxcore/subprocess"
	"jxcore/version"
	"jxcore/web"
	"jxcore/web/route"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

const (
	graceful = time.Second * 15
)

var (
	debug     bool   = false
	addr      string = ":80"
	noUpdate  bool   = false
	serverCmd        = &cobra.Command{
		Use:   "serve",
		Short: "Serve http backend for jxcore",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {

			Deamonize(func() {
				logger.Info("==================Jxcore Serve Starting=====================")
				currentdevice, err := device.GetDevice()
				if err != nil {
					logger.Fatal(err)
				}
				logger.Info("workerid : ", currentdevice.WorkerID)

				ctx, cancel := context.WithCancel(context.Background())
				errGroup, ctx := errgroup.WithContext(ctx)

				for dir, mapSrcDst := range monitor.GetMountCfg() {
					errGroup.Go(func() error { return monitor.MountListener(ctx, dir, mapSrcDst) })
				}

				ssdpClient := ssdp.NewClient(currentdevice.WorkerID, 5)
				errGroup.Go(func() error { return ssdpClient.Aliving(ctx) })

				if device.GetDeviceType() == version.Pro {
					logger.Info("=======================Configuring Network============================")
					core.ConfigNetwork()

					// Network interface auto switch
					// IoTEdge VPN auto reconnect
					errGroup.Go(func() error { return core.MaintainNetwork(ctx, noUpdate) })
				}

				logger.Info("================Configuring Environment===================")
				logger.Info("ensure tmpfs is mounted")
				err = ceph.EnsureTmpFs()
				if err != nil {
					logger.Fatal(err)
				}

				logger.Info("================Starting Subprocesses===================")
				subprocess.LoadConfig(viper.GetStringMap("components"))

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
					logger.WithFields(logger.Fields{"signal": sig}).Info("receive a signal to stop all process & exit")
					cancel()

					time.Sleep(graceful)
					logger.Info("===============Jxcore exited===============")
					logger.Fatal("Jxcore cannot exit within graceful period: ", graceful.String())
				}()

				err = errGroup.Wait()
				logger.Info("===============Jxcore exited===============")
				if err != nil {
					logger.Error("Exited with error: %+v", err)
				}
			})
		},
	}
)

func init() {
	serverCmd.PersistentFlags().StringVarP(&addr, "port", "p", ":80", "Addr to run Application server on")
	serverCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", true, "Whether to enable pprof")
	serverCmd.PersistentFlags().BoolVarP(&noUpdate, "no-update", "n", false, "Whether to check for update")
	serverCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "yaml setting for component")
	_ = viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("debug", serverCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("no-update", serverCmd.PersistentFlags().Lookup("no-update"))
	_ = viper.BindPFlag("config", serverCmd.PersistentFlags().Lookup("config"))

}
