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
	"jxcore/lowapi/utils"
	"jxcore/monitor"
	"jxcore/subprocess"
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
	debug    bool   = false
	addr     string = ":80"
	noUpdate bool   = false
	noDaemon bool   = false
	serveCmd        = &cobra.Command{
		Use:   "serve",
		Short: "Serve http backend for jxcore",
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
)

func init() {
	serveCmd.PersistentFlags().StringVarP(&addr, "port", "p", ":80", "Addr to run Application server on")
	serveCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", true, "Whether to enable pprof")
	serveCmd.PersistentFlags().BoolVarP(&noUpdate, "no-update", "n", false, "Whether to check for update")
	serveCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "yaml setting for component")
	serveCmd.PersistentFlags().BoolVar(&noDaemon, "no-daemon", noDaemon, "Debug mode: don't fork to the background")
	_ = viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("debug", serveCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("no-update", serveCmd.PersistentFlags().Lookup("no-update"))
	_ = viper.BindPFlag("config", serveCmd.PersistentFlags().Lookup("config"))

}

func serve() {
	logger.Info("==================Jxcore Serve Starting=====================")
	logger.Infof("Config: %+v", viper.GetViper().ConfigFileUsed())

	currentdevice, err := device.GetDevice()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("workerid : ", currentdevice.WorkerID)

	// 自动切换网卡
	// vpn 自动连接 IoTEdge
	// 保证 网络连接 是第一优先级，如果发生错误重启jxcore
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
	subprocess.LoadConfig(viper.GetStringMap("components"))

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
		logger.Info("===============Jxcore exited===============")
		logger.Fatal("Jxcore cannot exit within graceful period: ", graceful.String())
	}()

	_ = errGroup.Wait()

	logger.Info("===============Jxcore exited===============")
}
