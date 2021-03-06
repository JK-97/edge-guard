package register

import (
	"encoding/json"
	"github.com/JK-97/edge-guard/core/device"
	"github.com/JK-97/edge-guard/internal/config"
	"github.com/JK-97/edge-guard/internal/network/dns"
	"github.com/JK-97/edge-guard/internal/network/vpn"
	"github.com/JK-97/edge-guard/lowapi/logger"
	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/lowapi/utils"
	"io/ioutil"
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

	updateConsulConfig(currentdevice)

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
		err := ioutil.WriteFile(consulConfigPath, buf, 0666)
		if err == nil {
			go utils.RunAndLogError(func() error { return system.RunCommand("docker restart edgex-core-consul") })
		} else {
			logger.Error(err)
		}
	}
}

// updateTelegrafConfig 更新 Telegraf 和 InfluxDB ,cadvisor配置
func updateTelegrafConfig(currentdevice *device.Device, masterip string) {
	config.Telegrafcfg(masterip, currentdevice.WorkerID)
	config.Cadvisorcfg(masterip, currentdevice.WorkerID)
	VpnIP := vpn.GetClusterIP()
	if VpnIP != "" {
		config.Statsitecfg(masterip, VpnIP)
	}
	go utils.RunAndLogError(func() error {
		return system.RunCommand("docker-compose -f /jxbootstrap/worker/docker-compose.d/cadvisor/docker-compose.yaml down && docker-compose -f /jxbootstrap/worker/docker-compose.d/cadvisor/docker-compose.yaml up -d")
	})
	go utils.RunAndLogError(func() error {
		return system.RunCommand("docker-compose -f /jxbootstrap/worker/docker-compose.d/statsite/docker-compose.yml down && docker-compose -f /jxbootstrap/worker/docker-compose.d/statsite/docker-compose.yml up -d")
	})
}
