package dns

import (
	"io/ioutil"
	"strings"
)

// DNS方案：
// 建立虚拟网口lo:0，IP地址172.18.1.1作为dnsmasq的监听IP。
// /etc/resolv.conf 保持指向IP 172.18.1.1，使用dnsmasq
// dnsmasq 配置使用文件 /etc/dnsmasq.resolv.conf 作为上游DNS 地址
// edge-guard添加dhcp hook，使得dhcp获取的DNS写到/edge/resolv.d/dhclient.$interface
// 切换网卡时，将dhcp获取的DNS写入/etc/dnsmasq.resolv.conf

// ApplyInterfaceDNSResolv 将dhcp的resolv配置应用到dnsmasq
func ApplyInterfaceDNSResolv(iFace string) error {
	ifaceDNS, err := readInterfaceDnsResolv(iFace)
	if err != nil {
		return err
	}
	resolvData := readDNSResolvCustomContent() + ifaceDNSConfHeader + ifaceDNS
	return ioutil.WriteFile(dnsmasqResolvPath, []byte(resolvData), 0777)
}

func readDNSResolvCustomContent() (customContent string) {
	b, err := ioutil.ReadFile(dnsmasqResolvPath)
	if err != nil {
		return
	}
	l := strings.Split(string(b), "\n")
	for _, row := range l {
		if strings.HasPrefix(row, "#") && strings.Contains(row, "DHCP") {
			return
		}
		customContent += row + "\n"
	}
	return
}

func readInterfaceDnsResolv(iFace string) (string, error) {
	path := ifaceResolvPathPrefix + iFace
	ifaceData, err := ioutil.ReadFile(path)
	return string(ifaceData), err
}

func ParseInterfaceDNSResolv(iFace string) ([]string, error) {
	ifaceDNS, err := readInterfaceDnsResolv(iFace)
	if err != nil {
		return nil, err
	}

	var nameservers []string
	l := strings.Split(ifaceDNS, "\n")
	for _, row := range l {
		if strings.HasPrefix(row, "nameserver") {
			row = strings.TrimPrefix(row, "nameserver")
			row = strings.TrimSpace(row)
			nameservers = append(nameservers, row)
		}
	}
	return nameservers, nil
}
