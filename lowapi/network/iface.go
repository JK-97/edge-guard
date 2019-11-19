package network

import (
	"context"
	"fmt"
	"jxcore/config/yaml"
	"net"
	"os/exec"
	"time"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"

	"github.com/vishvananda/netlink"
)

const (
	highPriority = 5
)

var (
	currentIFace           string
	checkBestIFaceInterval = time.Second * 5
	ifacePriority          = yaml.Config.IFace.Priority
	backupIFace            = yaml.Config.IFace.Backup
)

func init() {
	interval, err := time.ParseDuration(yaml.Config.IFace.SwitchInterval)
	if err == nil {
		checkBestIFaceInterval = interval
	}
}

// 方案:
// 所有网卡连接上后，操作系统自动添加一条默认路由，metric（优先级）为100。
// jxcore选定的网卡，会添加metric为5的默认路由
// jxcore启动时调用InitIFace，选择优先级最高，能ping通外网的网卡
// 每隔checkBestIFaceInterval间隔，重新选择优先级最高，能ping通外网的网卡
func InitIFace() error {
	return switchIFace(findBestIFace())
}

func MaintainBestIFace() error {
	for {
		bestIFace := findBestIFace()
		err := switchIFace(bestIFace)
		if err != nil {
			log.Error("Failed to switch network interface: ", err)
		}
		time.Sleep(checkBestIFaceInterval)
	}
}

func findBestIFace() string {
	for _, iface := range ifacePriority {
		connected := testConnect(iface)
		if connected {
			return iface
		}
	}
	return backupIFace
}

func switchIFace(iFace string) (err error) {
	UnlockResolvConf()
	IFaceUp(iFace)
	LockResolvConf()
	route, err := getGWRoute(iFace)
	if err != nil {
		return err
	}
	route.Priority = highPriority
	err = netlink.RouteReplace(route)
	if err != nil {
		return err
	}

	currentIFace = iFace
	return
}

func testConnect(netInterface string) bool {
	if netInterface != currentIFace {
		IFaceUp(netInterface)
	}
	gwRoute, err := getGWRoute(netInterface)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); !ok {
			log.Info(err)
		}
		return false
	}
	dst := net.IPNet{
		IP:   net.ParseIP(testIP),
		Mask: net.CIDRMask(32, 32),
	}
	route := netlink.Route{
		Dst: &dst,
		Gw:  gwRoute.Gw,
	}

	err = netlink.RouteReplace(&route)
	if err != nil {
		log.Info(err)
		return false
	}
	return ping(testIP)
}

// getGWRoute 获取网卡的默认路由
func getGWRoute(netInterface string) (*netlink.Route, error) {
	link, err := netlink.LinkByName(netInterface)
	if err != nil {
		return nil, err
	}
	routes, _ := netlink.RouteList(link, netlink.FAMILY_V4)
	for _, r := range routes {
		if r.Gw != nil {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("Gateway route of %v not found", netInterface)
}

func IFaceUp(netInterface string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	_ = exec.CommandContext(ctx, "ifup", "--force", netInterface).Run()
	cancel()
}
