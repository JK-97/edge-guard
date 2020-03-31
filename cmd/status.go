package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JK-97/edge-guard/core/device"
	"github.com/JK-97/edge-guard/internal/network/vpn"

	log "github.com/JK-97/edge-guard/lowapi/logger"

	"github.com/spf13/cobra"
)

const (
	exitCodeNotInitialized int = 1 << iota
	exitCodeDHCPFailed
	exitCodeVPNFailed
)

// statusCmd 获取 edge-guard 的状态
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "to see the status of edge-guard",
	Long:  `to see the status of edge-guard`,

	Run: func(cmd *cobra.Command, args []string) {
		exitCode := 0
		info, err := device.GetDevice()
		if err != nil && os.IsNotExist(err) {
			log.Error("Not initialized.")
			exitCode |= exitCodeNotInitialized
			os.Exit(exitCode)
		}
		flags := cmd.PersistentFlags()
		if ok, _ := flags.GetBool("device"); ok {
			log.Info("WorkID: ", info.WorkerID)
			log.Info("DhcpServer: ", info.DhcpServer)
			log.Info("DeviceKey: ", info.Key)
			log.Info("VPN Mode: ", info.Vpn)
		}

		fmt.Println("Connect to DHCP Server")
		if resp, err := http.Get(info.DhcpServer); err != nil {
			log.Error(err)
			exitCode |= exitCodeDHCPFailed
			log.Error("Connect Failed.")
		} else if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
			fmt.Println(resp.StatusCode)
			exitCode |= exitCodeDHCPFailed
			log.Error("Connect Failed.")
		} else {
			log.Info("Connect Success.")
		}

		if ok, _ := flags.GetBool("vpn"); ok && info.Vpn != device.VPNModeLocal {
			log.Info("Test VPN Status")
			ip := vpn.GetClusterIP()
			if ip != "" {
				log.Info("VPN Test Success, ClusterIP: ", ip)
			} else {
				exitCode |= exitCodeVPNFailed
				log.Error("VPN Test Failed!")
			}
		}

		os.Exit(exitCode)
	}}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.AddCommand(statusCmd)

	flags := statusCmd.PersistentFlags()
	flags.BoolP("device", "d", true, "Print device informations.")
	flags.BoolP("vpn", "v", true, "Test VPN Status")
	flags.BoolP("gateway", "g", true, "Test Gateway")

}
