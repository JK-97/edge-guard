package cmd

import (
	"fmt"
	"jxcore/core/device"
	"jxcore/internal/network/vpn"
	"net/http"
	"os"

	log "jxcore/lowapi/logger"

	"github.com/spf13/cobra"
)

const (
	exitCodeNotInitialized int = 1 << iota
	exitCodeDHCPFailed
	exitCodeVPNFailed
)

// statusCmd 获取 jxcore 的状态
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "to see the status of jxcore",
	Long:  `to see the status of jxcore`,

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
	flags := statusCmd.PersistentFlags()
	flags.BoolP("device", "d", true, "Print device informations.")
	flags.BoolP("vpn", "v", true, "Test VPN Status")
	flags.BoolP("gateway", "g", true, "Test Gateway")
}
