package core

import (
	"context"
	"jxcore/config/yaml"
	"jxcore/core/device"
	"jxcore/core/hearbeat"
	"jxcore/core/register"
	"jxcore/internal/network"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/iface"
	"jxcore/lowapi/docker"
	"jxcore/management/updatemanage"
	"time"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"

	"golang.org/x/sync/errgroup"
)

func NewJxCore() *JxCore {
	return &JxCore{}
}

func GetJxCore() *JxCore {
	return jxcore
}

func (j *JxCore) ConfigSupervisor() {
	startupProgram := yaml.Config
	yaml.ParseAndCheck(*startupProgram, "")
}

func (j *JxCore) ConfigNetwork() {
	// NetworkManager 和 systemd-resolved 会更改 /etc/resolv.conf，使dns不可控。需要停止
	network.DisableNetworkManager()
	network.DisableSystemdResolved()

	// nano: /etc/dhcp/dhclient-enter-hooks.d 里的脚本需要被移除。
	// 如果脚本存在，`ifup`时不会更新 /etc/resolv.conf
	dns.RemoveDHCPEnterHooks()

	// 将dhcp 获取的resolv 信息重定向到 /edge/resolv.d
	dns.ApplyDHCPResolveUpdateHooks()
	dns.ResetResolv()
	dns.ResetDNSMasqConf()

	err := iface.InitIFace()
	if err != nil {
		log.Error(err)
	}
	dns.ResetHostFile(network.GetEthIP())

	err = docker.EnsureDockerDNSConfig()
	if err != nil {
		log.Error("Failed to configure docker DNS: ", err)
	}
}

func (j *JxCore) MaintainNetwork(ctx context.Context) error {
	errGroup := errgroup.Group{}
	errGroup.Go(iface.MaintainBestIFace)
	errGroup.Go(func() error { return maintainMasterConnection(ctx) })
	return errGroup.Wait()
}

//contrl the update
func (j JxCore) UpdateCore() {
	for !network.CheckMasterConnect() {
		time.Sleep(5 * time.Second)
		log.Info("Waiting for master connect")

	}
	log.Info("Master Connect")
	if dns.CheckDnsmasqConf() {
		log.Info("Normal Dnsmasq configuration ")
	} else {
		log.Error("Error Dnsmasq configuration ")
	}
	updatemanage.AddAptKey()
	updateprocess := updatemanage.GetUpdateProcess()
	pkgneedupdate := updateprocess.CheckUpdate()
	if len(pkgneedupdate) != 0 {
		updateprocess.UpdateSource()
		updateprocess.UpdateComponent(pkgneedupdate)
	}
	updateprocess.ReportVersion()
}

func maintainMasterConnection(ctx context.Context) error {
	var mymasterip string
	currentedvice, err := device.GetDevice()
	if err != nil {
		return err
	}
	for {
		for {
			register.FindMasterFromDHCPServer(currentedvice.WorkerID, currentedvice.Key)
			//获取vpn key，连接vpn
			mymasterip, err = register.GetMyMaster(currentedvice.WorkerID, currentedvice.Key)
			//校验新的master是否协力hossts文件
			time.Sleep(3 * time.Second)
			log.Info("Register Worker Net", mymasterip)
			if err == nil {
				break
			}

		}
		time.Sleep(3 * time.Second)

		// VPN 就绪之后 启动 component 按照配置启动(同步工具集合)
		hearbeat.AliveReport(mymasterip)
	}
}
