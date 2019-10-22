package core

import (
	"bytes"
	"net"
	"os/exec"
	"strconv"
	"text/template"

	log "jxcore/go-utils/logger"

	"github.com/vishvananda/netlink"
)

// 推荐网卡名/次要网卡名
var (
	RecommendIface = "eth0"
	SecondaryIface = "usb0"

	useEthTpl *template.Template
	useUSBTpl *template.Template

	usbNetworkReachable = false // USB 网卡不可用
)

func main() {
	SetUp()

	if usbNetworkReachable {
		done := make(chan struct{})
		go linkSubscribe(done)

		<-done
	}
}

// SetUp 初始化设置
func SetUp() {

	useEthTpl = template.Must(template.New("eth").Parse(`ip route del default; dhclient {{.Iface}}; killall dhclient`))
	useUSBTpl = template.Must(template.New("usb").Parse(`ip route del default; ip route add default dev {{.Iface}}`))

	checkUSBEnable()
}

const (
	testServer = "114.114.114.114"
	testPort   = 53
)

// 检测 USB 网卡是否可用
func checkUSBEnable() {

	if _, err := netlink.LinkByName(RecommendIface); err == nil {
		cmd := exec.Command("dhclient", RecommendIface)
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
	}

	_, err := netlink.LinkByName(SecondaryIface)
	if err != nil {
		log.Info(SecondaryIface, " not ready.")
		return
	}

	enableUSB()

	if tcping(testServer, testPort) {
		usbNetworkReachable = true
	} else {
		// if netlink.LinkByName()
	}
	enableEthernet()

	// 有线网卡不通
	if !tcping(testServer, testPort) {
		enableUSB()
	}
}

type iface struct {
	Iface string
}

// enableEthernet 启用以太网 作为默认路由
func enableEthernet() error {
	buffer := bytes.NewBuffer(nil)
	useEthTpl.Execute(buffer, iface{RecommendIface})
	shellScript := buffer.String()
	cmd := exec.Command("sh", "-c", shellScript)
	err := cmd.Run()

	return err
}

// enableUSB 启用 USB 网卡 作为默认路由
func enableUSB() error {
	buffer := bytes.NewBuffer(nil)
	useUSBTpl.Execute(buffer, iface{SecondaryIface})
	shellScript := buffer.String()
	cmd := exec.Command("sh", "-c", shellScript)
	err := cmd.Run()
	// TODO: 修改 DNS 服务器
	return err
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
			// log.Info("Name: ", attrs.Name)

			switch attrs.OperState {
			case netlink.OperUp:
				log.Info("Name: ", attrs.Name, " Up")
				if attrs.Name == RecommendIface {
					enableEthernet()
					// TODO 测试 网络是否可达
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
			// ip route del default; ip route add default dev eth0 via 10.55.2.253
			// ip route del default; ip route add default dev usb0
		case addr := <-addrChan:
			if len(addr.LinkAddress.IP) == net.IPv4len {
				if addr.NewAddr {
					log.Info("Up ", addr)
				} else {
					log.Info("Down ", addr)
				}
			} else if len(addr.LinkAddress.IP) == net.IPv6len {
				log.Debug("IPv6")
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
		log.Info("Error", err)
		return false
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	err = tcpConn.CloseWrite()
	if err != nil {
		log.Info("Error", err)
		return false
	}

	return true
}
