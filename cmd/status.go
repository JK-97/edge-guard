package cmd

import (
	"fmt"
	"jxcore/regeister"
	"os"

	"github.com/spf13/cobra"
)

const (
	exitCodeNotInitialized int = 1 << iota
	exitCodeVPNFailed
)

// statusCmd 获取 jxcore 的状态
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "to see the status of jxcore",
	Long:  `to see the status of jxcore`,

	Run: func(cmd *cobra.Command, args []string) {
		exitCode := 0
		info, err := regeister.ReadDeviceInfo()
		if err != nil && os.IsNotExist(err) {
			fmt.Println("Not initialized.")
			exitCode |= exitCodeNotInitialized
			os.Exit(exitCode)
		}
		flags := cmd.PersistentFlags()
		if ok, _ := flags.GetBool("device"); ok {
			fmt.Println("WorkID:", info.WorkID)
			fmt.Println("DhcpServer:", info.DhcpServer)
			fmt.Println("DeviceKey:", info.Key)
			fmt.Println("VPN Mode:", info.Vpn)
		}
		if ok, _ := flags.GetBool("vpn"); ok && info.Vpn != regeister.VPNModeLocal {
			fmt.Println("Test VPN Status")
			var ip string
			switch info.Vpn {
			case regeister.VPNModeOPENVPN:
				regeister.Closeopenvpn()
				regeister.Startopenvpn()
			case regeister.VPNModeWG:
				regeister.CloseWg()
				regeister.StartWg()
			}
			ip = regeister.GetClusterIP()
			if ip != "" {
				fmt.Println("ClusterIP:", ip)
			} else {
				exitCode |= exitCodeVPNFailed
				fmt.Println("VPN Test Failed!")
			}
		}

		// if ok, _ := flags.GetBool("gateway"); ok {
		// 	fmt.Println("Test Start Gateway")
		// 	monitor.GWEmitter()
		// }

		os.Exit(exitCode)
	}}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.AddCommand(statusCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flags := statusCmd.PersistentFlags()
	flags.BoolP("device", "d", true, "Print device informations.")
	flags.BoolP("vpn", "v", true, "Test VPN Status")
	flags.BoolP("gateway", "g", true, "Test Gateway")
	// flags.BoolP("mongo", "m", false, "Recover Mongo")
}
