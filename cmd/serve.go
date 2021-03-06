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

	"github.com/JK-97/edge-guard/config"
	"github.com/JK-97/edge-guard/config/yaml"
	"github.com/JK-97/edge-guard/core"
	"github.com/JK-97/edge-guard/core/device"
	"github.com/JK-97/edge-guard/gateway"
	"github.com/JK-97/edge-guard/internal/network/ssdp"
	"github.com/JK-97/edge-guard/lowapi/ceph"
	"github.com/JK-97/edge-guard/lowapi/logger"
	"github.com/JK-97/edge-guard/lowapi/utils"
	"github.com/JK-97/edge-guard/monitor"
	"github.com/JK-97/edge-guard/subprocess"
	"github.com/JK-97/edge-guard/web"
	"github.com/JK-97/edge-guard/web/route"

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
	noDaemon bool   = false
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve http backend for edge-guard",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if noDaemon {
			serve()
		} else {
			Deamonize(serve)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&addr, "port", addr, "Addr to run Application server on")
	serveCmd.PersistentFlags().BoolVar(&debug, "debug", debug, "Whether to enable pprof")
	serveCmd.PersistentFlags().BoolVar(&noUpdate, "no-update", noUpdate, "Whether to check for update")
	serveCmd.PersistentFlags().BoolVar(&noDaemon, "no-daemon", noDaemon, "Debug mode: don't fork to the background")

	serveCmd.PersistentFlags().String("interface", "eth0", "gateway listen where")
	serveCmd.PersistentFlags().String("config", "./settings.yaml", "yaml setting for component")
	cfg := config.Config()
	_ = cfg.BindPFlag("yamlsettings", serveCmd.PersistentFlags().Lookup("config"))
	_ = cfg.BindPFlag("interface", serveCmd.PersistentFlags().Lookup("interface"))
}

func serve() {
	logger.Info("==================edge-guard Serve Starting=====================")
	logger.Infof("Config: %+v", yaml.Config)

	currentdevice, err := device.GetDevice()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("workerid : ", currentdevice.WorkerID)

	// 自动切换网卡
	// vpn 自动连接 IoTEdge
	// 保证 网络连接 是第一优先级，如果发生错误重启edge-guard
	logger.Info("=======================Configuring Network============================")
	ctx, cancel := context.WithCancel(context.Background())
	errGroup, ctx := errgroup.WithContext(ctx)
	go func() {
		err := core.MaintainNetwork(ctx, noUpdate)
		logger.Error("Maintain network error: ", err)
		cancel()
	}()

	logger.Info("================Configuring File System===================")
	logger.Info("ensure tmpfs is mounted")
	go utils.RunAndLogError(ceph.EnsureTmpFs)
	logger.Info("init sdcard mount")
	monitor.InitTF()
	logger.Info("start auto sdcard mount")
	for dir, mapSrcDst := range monitor.GetMountCfg() {
		utils.GoAndRestartOnError(ctx, errGroup, "mount listener "+dir, func() error { return monitor.MountListener(ctx, dir, mapSrcDst) })
	}

	logger.Info("================Starting Subprocesses===================")
	core.ConfigSupervisor()

	// start ssdp
	ssdpClient := ssdp.NewClient(currentdevice.WorkerID, 5)
	utils.GoAndRestartOnError(ctx, errGroup, "ssdp", func() error { return ssdpClient.Aliving(ctx) })

	// start up all component process
	utils.GoAndRestartOnError(ctx, errGroup, "subprocess", func() error { return subprocess.RunServer(ctx) })

	// web server
	utils.GoAndRestartOnError(ctx, errGroup, "gateway", func() error { return gateway.ServeGateway(ctx, graceful) })
	utils.GoAndRestartOnError(ctx, errGroup, "web", func() error { return web.Serve(ctx, addr, route.Routes(), graceful) })
	utils.GoAndRestartOnError(ctx, errGroup, "prof web", func() error { return web.Serve(ctx, ":10880", http.DefaultServeMux, graceful) })

	// 处理 SIGTERM 和 SIGINT
	go func() {
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-termChan
		logger.WithFields(logger.Fields{"signal": sig}).Info("receive a signal to stop all process & exit")
		cancel()
	}()

	<-ctx.Done()
	// fatal after graceful period
	go func() {
		time.Sleep(graceful)
		logger.Info("===============edge-guard exited===============")
		logger.Fatal("edge-guard cannot exit within graceful period: ", graceful.String())
	}()

	_ = errGroup.Wait()

	logger.Info("===============edge-guard exited===============")
}
