package iface

import (
	"context"
	"fmt"
	"jxcore/config/yaml"
	"jxcore/core/device"
	"jxcore/gateway/log"
	"jxcore/internal/network"
	"jxcore/internal/network/dns"
	"net"
	"net/url"
	"os/exec"
	"time"

	"github.com/vishvananda/netlink"
)

const (
	testIP       = "114.114.114.114"
	highPriority = 5
	dhcpPriority = 6
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

// 网卡切换方案:
// 所有网卡连接上后，操作系统自动添加一条默认路由，metric（优先级）为100+。
// jxcore选定的网卡，会添加metric为5的默认路由
// jxcore启动时调用InitIFace，选择优先级最高，能ping通外网的网卡
// 每隔checkBestIFaceInterval间隔，重新选择优先级最高，能ping通外网的网卡

func MaintainBestIFace(ctx context.Context) error {
	// parse the dhcpserver ip
	deviceInfo, err := device.GetDevice()
	if err != nil {
		return err
	}
	urlInfo, err := url.Parse(deviceInfo.DhcpServer)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(checkBestIFaceInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			bestIFace := findBestIFace()
			err = switchIFace(bestIFace, urlInfo.Hostname())
			if err != nil {
				log.Error("Failed to switch network interface: ", err)
			}
		}
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

func switchIFace(iface, dhcpServer string) (err error) {
	if iface == currentIFace {
		return nil
	}
	route, err := GetGWRoute(iface)
	if err != nil {
		return err
	}
	servers, err := net.LookupHost(dhcpServer)
	if err != nil {
		return err
	}
	// dhcpserver maybe not only one
	for _, addr := range servers {
		route.Priority = dhcpPriority
		route.Dst = &net.IPNet{
			// only add route to dhcpserver,do not change system default route
			IP:   net.ParseIP(addr),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}

		err = netlink.RouteReplace(route)
		if err != nil {
			return err
		}
	}

	err = dns.ApplyInterfaceDNSResolv(iface)
	if err != nil {
		IFaceUp(iface)
		err = dns.ApplyInterfaceDNSResolv(iface)
		if err != nil {
			return err
		}
	}

	log.Infof("Switch network interface %s -> %s", currentIFace, iface)
	currentIFace = iface
	return
}

func testConnect(netInterface string) bool {
	gwRoute, err := GetGWRoute(netInterface)
	if err != nil {
		// route not exists, ifup then retry
		IFaceUp(netInterface)
		gwRoute, err = GetGWRoute(netInterface)
	}
	if err != nil {
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
	if network.Ping(testIP) {
		return true
	}
	IFaceUp(netInterface)
	return network.Ping(testIP)
}

// GetGWRoute 获取网卡的默认路由
func GetGWRoute(netInterface string) (*netlink.Route, error) {
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

// 刷新网卡配置，自动添加 route, 添加 /edge/resolv.d/dhclient.$interface
func IFaceUp(netInterface string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	_ = exec.CommandContext(ctx, "ifup", "--force", netInterface).Run()
	cancel()
	_ = exec.Command("killall", "dhclient").Run()
}

func GetCurrentIFcae() string {
	return currentIFace
}

func SetHighPriority(route *netlink.Route) {
	route.Priority = highPriority
}
