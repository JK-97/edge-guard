package register

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/internal/network"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/vpn"
	"jxcore/internal/template"
	"jxcore/lowapi/utils"
)

type consulConfig struct {
	Server           bool     `json:"server"`
	ClientAddr       string   `json:"client_addr"`
	AdvertiseAddrWan string   `json:"advertise_addr_wan"`
	BootstrapExpect  int      `json:"bootstrap_expect"`
	Datacenter       string   `json:"datacenter"`
	NodeName         string   `json:"node_name"`
	RetryJoinWan     []string `json:"retry_join_wan"`
	UI               bool     `json:"ui"`
}

// onMasterIPChanged master IP 变化后执行
func onMasterIPChanged(masterip string) {
	currentdevice, err := device.GetDevice()
	utils.CheckErr(err)

	if utils.FileExists(consulConfigPath) {
		updateConsulConfig(currentdevice)
	}

	updateTelegrafConfig(currentdevice, masterip)
}

// updateConsulConfig 更新 Consul 配置
func updateConsulConfig(currentdevice *device.Device) {
	config := consulConfig{
		Server:           true,
		ClientAddr:       "0.0.0.0",
		AdvertiseAddrWan: vpn.GetClusterIP(),
		BootstrapExpect:  1,
		Datacenter:       "worker-" + currentdevice.WorkerID,
		NodeName:         "worker-" + currentdevice.WorkerID,
		RetryJoinWan:     []string{dns.MasterHostName},
		UI:               true,
	}
	if buf, err := json.Marshal(config); err == nil {
		ioutil.WriteFile(consulConfigPath, buf, 0666)
	}
}

// updateTelegrafConfig 更新 Telegraf 和 InfluxDB ,cadvisor配置
func updateTelegrafConfig(currentdevice *device.Device, masterip string) {
	template.Telegrafcfg(masterip, currentdevice.WorkerID)
	template.Cadvisorcfg(masterip, currentdevice.WorkerID)
	var VpnIP string
	//确保4g 或 以太有一个起来的情况下
	if _, erreth0 := network.GetMyIP("eth0"); erreth0 == nil {
		VpnIP = vpn.GetClusterIP()
	} else if _, errusb0 := network.GetMyIP("usb0"); errusb0 == nil {
		VpnIP = vpn.GetClusterIP()
	}
	if VpnIP != "" {
		template.Statsitecfg(masterip, VpnIP)
	}
}
