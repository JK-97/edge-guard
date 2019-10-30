package core

import (
	"bufio"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "jxcore/go-utils/logger"
	"jxcore/lowapi/dns"

	"github.com/vishvananda/netlink"
)

// 推荐网卡名/次要网卡名
var (
	RecommendIface = "eth0"
	SecondaryIface = "usb0"

	usbNetworkReachable = false // USB 网卡不可用

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
				log.Info("Name: ", attrs.Name, " Up")
				if attrs.Name == RecommendIface {
					enableEthernet()
					// 有线网卡不通
					if !tcping(testServer, testPort) {
						enableUSB()
					}
				}
			case netlink.OperDown:
				log.Info("Name: ", attrs.Name, " Down")
				if attrs.Name == RecommendIface {
					enableUSB()
				}
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
