package cmd

import (
    "fmt"
    "jxcore/core/device"

    "jxcore/lowapi/network"

    "jxcore/lowapi/vpn"
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
        cuerrentdevice, err := device.GetDevice()
        if err != nil && os.IsNotExist(err) {
            fmt.Println("Not initialized.")
            exitCode |= exitCodeNotInitialized
            os.Exit(exitCode)
        }
        flags := cmd.PersistentFlags()
        if ok, _ := flags.GetBool("device"); ok {
            fmt.Println("WorkerID:", cuerrentdevice.WorkerID)
            fmt.Println("DhcpServer:", cuerrentdevice.DhcpServer)
            fmt.Println("DeviceKey:", cuerrentdevice.Key)
            fmt.Println("VPN Mode:", cuerrentdevice.Vpn)
        }
        if ok, _ := flags.GetBool("vpn"); ok && cuerrentdevice.Vpn != device.VPNModeLocal {
            fmt.Println("Test VPN Status")
            var ip string
            switch cuerrentdevice.Vpn {
            case device.VPNModeOPENVPN:
                vpn.Closeopenvpn()
                vpn.Startopenvpn()
            case device.VPNModeWG:
                vpn.CloseWg()
                vpn.StartWg()
            }
            ip = network.GetClusterIP()
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

    flags := statusCmd.PersistentFlags()
    flags.BoolP("device", "d", true, "Print device informations.")
    flags.BoolP("vpn", "v", true, "Test VPN Status")
    flags.BoolP("gateway", "g", true, "Test Gateway")

}
