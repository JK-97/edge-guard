package core

import (
	"bufio"
	"errors"
	"io/ioutil"
	"jxcore/config/yaml"
	"jxcore/core/device"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "jxcore/go-utils/logger"
	"jxcore/lowapi/dns"

	"github.com/rjeczalik/notify"
	"github.com/vishvananda/netlink"
)

// 推荐网卡名/次要网卡名
var (
	routeEchoFile  = "/edge/route"
	RecommendIface = "eth0"
	SecondaryIface = "usb0"

	usbNetworkReachable = false // USB 网卡不可用
	debug               = false
)

func main() {
	SetUp()

	// if usbNetworkReachable {
	// }
	done := make(chan struct{})
	go linkSubscribe(done)

	<-done
}

const leaseFilename = "/tmp/dhclient.leases"

// SetUp 初始化设置
func SetUp() {

	checkUSBEnable()
}

const (
	testServer = "114.114.114.114"
	testPort   = 53
)

type leasesConfig struct {
	DomainNameServer string
	Router           string
	// fixed-address 192.168.225.23;
	// option subnet-mask 255.255.255.0;
	// option dhcp-lease-time 43200;
	// option routers 192.168.225.1;
	// option dhcp-message-type 5;
	// option dhcp-server-identifier 192.168.225.1;
	// option domain-name-servers 192.168.225.1;
	// option interface-mtu 1500;
	// option dhcp-renewal-time 21600;
	// option unknown-120 0:4:61:62:63:64:3:63:6f:6d:0;
	// option dhcp-rebinding-time 37800;
	// option broadcast-address 192.168.225.255;
	// option host-name "worker-J013fcec6f";
}

// 解析 dhclient.leases 文件
func parseLeaseFile(filename string) *leasesConfig {
	if filename == "" {
		filename = leaseFilename
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer f.Close()

	conf := &leasesConfig{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "routers") {
			line = strings.TrimSpace(line)
			segs := strings.Split(line, "routers")
			if len(segs) > 1 {
				segs = strings.Split(strings.TrimSpace(segs[1]), " ")
				conf.Router = strings.Trim(segs[0], ";")
			}
		} else if strings.Contains(line, "domain-name-servers") {
			line = strings.TrimSpace(line)
			segs := strings.Split(line, "domain-name-servers")
			if len(segs) > 1 {
				segs = strings.Split(strings.TrimSpace(segs[1]), " ")
				conf.DomainNameServer = strings.Trim(segs[0], ";")
			}
		}
	}
	log.Debug("dhclient", conf)

	return conf
}

// 检测 USB 网卡是否可用
func checkUSBEnable() {
	log.Debug("checkUSBEnable")

	if _, err := netlink.LinkByName(RecommendIface); err == nil {
		enableEthernet()
	}

	_, err := netlink.LinkByName(SecondaryIface)
	if err != nil {
		log.Info(SecondaryIface, " not ready.")
		return
	}

	enableUSB()

	// time.Sleep(time.Millisecond * 500)

	if tcping(testServer, testPort) {
		usbNetworkReachable = true
	}

	if err := enableEthernet(); err != nil {
		log.Error("checkUSBEnable: ", err)
	}

	// 有线网卡不通
	if !tcping(testServer, testPort) {
		enableUSB()
	}
}

type iface struct {
	Iface string
}

// 启用指定的网卡
func enableNetworkInterface(iface string) error {
	logger := log.WithFields(log.Fields{"chnw": iface})
	os.Remove(leaseFilename)

	cmd := exec.Command("dhclient", iface, "-lf", leaseFilename)
	log.Info(cmd.Args)
	err := cmd.Run()

	if err != nil {
		return err
	}

	conf := parseLeaseFile(leaseFilename)
	if conf != nil && conf.Router != "" {
		cmd = exec.Command("ip", "route", "replace", "default", "via", conf.Router, "dev", iface)

		logger.Info(cmd.Args)
		err = cmd.Run()
		if err == nil {
			dns.ResolvGuard()
		}
	} else {
		logger.Warn("Parse Failed.")
	}

	return err
}

// enableEthernet 启用以太网 作为默认路由
func enableEthernet() error {
	return enableNetworkInterface(RecommendIface)
}

// enableUSB 启用 USB 网卡 作为默认路由
func enableUSB() error {
	return enableNetworkInterface(SecondaryIface)
}

func linkSubscribe(done <-chan struct{}) {
	linkChan := make(chan netlink.LinkUpdate)
	// done := make(chan struct{})

	log.Info("Begin to check link")
	err := netlink.LinkSubscribe(linkChan, done)
	if err != nil {
		log.Fatal(err)
	}

	addrChan := make(chan netlink.AddrUpdate)
	log.Info("Begin to check addr")
	err = netlink.AddrSubscribe(addrChan, done)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case link := <-linkChan:
			// log.Println("IfInfomsg: ", link.IfInfomsg)
			// log.Println("Header: ", link.Header)
			// log.Println("Link: ", link.Link)
			if link.Link == nil {
				continue
			}
			attrs := link.Link.Attrs()

			if strings.HasPrefix(attrs.Name, "veth") {
				continue
			}

			switch attrs.OperState {
			case netlink.OperUp:
				routeHandler()

			case netlink.OperDown:
				// initRoute()
			default:
				log.Info("Name: ", attrs.Name, " OperState: ", attrs.OperState)
			}
			// ip link set usb0 down
			// ip link set usb0 up
			// ip route replace default dev eth0 via 10.55.2.253
			// ip route replace default dev usb0
		case addr := <-addrChan:
			if len(addr.LinkAddress.IP) == net.IPv4len {
				if addr.NewAddr {
					log.Info("Up ", addr)
				} else {
					log.Info("Down ", addr)
				}
			} else if len(addr.LinkAddress.IP) == net.IPv6len {
				// log.Debug("IPv6")
			}
		}
	}
}

func tcping(s string, port int) bool {
	if port <= 0 {
		port = 80
	}

	portStr := strconv.Itoa(port)

	addr := net.JoinHostPort(s, portStr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Info("Error: ", err)
		return false
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	err = tcpConn.CloseWrite()
	if err != nil {
		log.Info("Error: ", err)
		return false
	}

	return true
}

// RouteDetector 检测 rouote 文件的改动,并设置路由
func routeDetector() {
	c := make(chan notify.EventInfo, 2)
	if err := notify.Watch(routeEchoFile, c, notify.All); err != nil {
		log.Error(err)
	}
	for ei := range c {
		switch ei.Event() {
		case notify.Remove:
			go routeHandler()
			notify.Stop(c)
			routeDetector()
		}

	}
}

func RunRouteDetector() {
	if err := initRoute(); err != nil {
		log.Error("no Available interface cant be use")
	}
	log.Info("finished init route")
	go routeDetector()

}

type route struct {
	router string
	iface  string
}

var ifacePriority = map[string]int{"eth0": 1, "usb0": 2, "usb1": 3}

// routeHandler 处理函数, 自动选择优先级最高
func routeHandler() {
	// 根据优先级 配置 路由
	toChangeRoute, err := getRouterFromFile()
	if err != nil {
		return
	}

	if defaultRouter, err := getDefaultRouter(); err == nil {
		log.Info("wantToUseRoute: ", toChangeRoute.iface)
		log.Info("currentDefaultRoute: ", defaultRouter.iface)

		if isHigherThan(toChangeRoute.iface, defaultRouter.iface) {
			setDefaultRoute(toChangeRoute.iface, toChangeRoute.router)
		} else {
			log.Info("The currently used interface has a higher priority and does not change.")
		}

	} else {
		setDefaultRoute(toChangeRoute.iface, toChangeRoute.router)
	}

	// finish <- true
}

// setDefaultRoute 设置默认路由
func setDefaultRoute(iface, router string) (err error) {
	log.Info("use " + router + "--" + iface)
	err = exec.Command("ip", "route", "replace", "default", "via", router, "dev", iface).Run()
	waitForSetDefault(iface)
	return err
}
func getDhcpHost() (dhcpHost string) {
	currentdevice, _ := device.GetDevice()

	//解析dhcp url ， 获取 host
	dhcpServerInfo, err := url.Parse(currentdevice.DhcpServer)
	if err != nil {
		log.Info(err)
		return
	}
	dhcpHost, _, _ = net.SplitHostPort(dhcpServerInfo.Host)
	if debug {
		dhcpHost = "114.114.114.114"
	}

	return
}

// getAbleRouter 获取所有可用的 interface route
func getRouterFromFile() (changeRoute route, err error) {

	if _, err = os.Stat(routeEchoFile); err != nil {
		log.Error(err)
		return
	}
	dataByte, err := ioutil.ReadFile(routeEchoFile)
	if err != nil {
		log.Error(err)
		return
	}
	routeLine := strings.Split(strings.TrimSpace(string(dataByte)), " ")
	// 10.55.2.253 eth0
	netRoute := strings.TrimSpace(routeLine[0])
	netInterface := strings.TrimSpace(routeLine[1])

	changeRoute = route{router: netRoute,
		iface: netInterface}
	return
}

func setIPRoute(netRoute, netInterface string) (err error) {
	// err = exec.Command("ip", "route", "add", "114.114.114.114/32", "via", netRoute, "dev", netInterface).Run()
	err = exec.Command("/bin/bash", "-c", "ip route replace 114.114.114.114/32 via "+netRoute+" dev "+netInterface).Run()
	log.Info("setIPRoute:", netInterface)
	waitForSet(netInterface)
	return
}

func removeIPRoute(netRoute, netInterface string) (err error) {
	err = exec.Command("/bin/bash", "-c", "ip route del 114.114.114.114/32 via "+netRoute+" dev "+netInterface).Run()
	log.Info("removeIPRoute :", netInterface)
	waitForRemove(netInterface)
	return
}

func checkRoute(netRoute, netInterface, dhcpHost string) (err error) {
	setIPRoute(netRoute, netInterface)
	err = exec.Command("ping", "-c", "1", "-I", netInterface, dhcpHost).Run()

	if err != nil {
		// 还在 调试中，ping -c时以code2 退出
		if err.Error() == "exit status 2" {
			err = nil
		} else {
			log.Info("checkRoute Unavailable : ", netInterface)
		}

	}

	// time.Sleep(30 * time.Second)
	removeIPRoute(netRoute, netInterface)

	return
}

func getDefaultRouter() (defaultRoute route, err error) {
	output, _ := exec.Command("ip", "route").Output()
	currentRouteInfo := strings.Split(string(output), "\n")
	if strings.Contains(currentRouteInfo[0], "default") {
		tmp := strings.Split(currentRouteInfo[0], " ")
		netInterface := tmp[len(tmp)-2]
		netRouter := tmp[len(tmp)-4]
		defaultRoute = route{router: netRouter, iface: netInterface}
		return
	} else {
		err = errors.New("can not find default interface")
		return
	}
}

func sortRouteByPriority(routeList []route) (recommendRoute route) {
	//比较优先级别,
	var min = ifacePriority[routeList[0].iface]
	for _, route := range routeList {
		if ifacePriority[route.iface] <= min {
			min = ifacePriority[route.iface]
			recommendRoute = route
			log.Info(recommendRoute.router, recommendRoute.iface)
		}
	}

	return
}

func isHigherThan(primary, secondary string) bool {
	return ifacePriority[primary] <= ifacePriority[secondary]
}

func initRoute() (err error) {
	dhcpHost := getDhcpHost()
	routerList := []route{}
	for iface, _ := range ifacePriority {

		route, err := LookUpIface(iface, dhcpHost)
		log.Info("find route : ", iface, route.router, route.iface)
		if err != nil {
			log.Info(err)
		} else {
			routerList = append(routerList, route)
		}

	}
	log.Info(routerList)
	if len(routerList) != 0 {
		recommendRoute := sortRouteByPriority(routerList)
		log.Info("recommendRoute", recommendRoute, len(routerList))
		err = setDefaultRoute(recommendRoute.iface, recommendRoute.router)
	} else {
		err = errors.New("no route can find")
	}
	return
}

func LookUpIface(iface, dhcpHost string) (theRoute route, err error) {
	theRoute = route{}
	exec.Command("dhclient", iface).Run()
	time.Sleep(2 * time.Second)

	theRoute, err = getRouterFromFile()
	if err != nil {
		log.Info(err)
		return
	}

	if theRoute.iface != iface {
		err = errors.New("cant not find " + iface)
		return
	}

	// err = checkRoute(theRoute.router, theRoute.iface, dhcpHost)
	checkRoute(theRoute.router, theRoute.iface, dhcpHost)

	if err != nil {
		log.Info(err)
		return
	}
	log.Info(err)
	return
}

func waitForSet(netInterface string) {

	for {
		output, err := exec.Command("ip", "route").Output()
		if err != nil {
			log.Info(err)
			return
		}

		lines := strings.Split(string(output), "\n")

		for _, line := range lines {
			if strings.Contains(line, "114.114.114.114") || strings.Contains(line, netInterface) {
				log.Info(line)
				goto endfor
			}
		}
	endfor:
		break
	}

}

func waitForRemove(netInterface string) {
	for {
		output, err := exec.Command("ip", "route").Output()
		if err != nil {
			log.Info(err)
			return
		}

		if !strings.Contains(string(output), "114.114.114.114") {
			break
		}
	}

}
func waitForSetDefault(netInterface string) {
	for {
		output, err := exec.Command("ip", "route").Output()
		if err != nil {
			log.Info(err)
			return
		}

		if !strings.Contains(string(output), "defalut") {
			break
		}
	}

}
func init() {
	if settings, err := yaml.LoadYaml("/edge/jxcore/bin/settings.yaml"); err == nil {
		debug = settings.Debug
		log.Info("debug model ", debug)
	}

}
