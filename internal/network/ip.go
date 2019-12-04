package network

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

func GetMyIP(name string) (string, error) {
	ips, err := GetMyIPSlice(name)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("No IP found for interface %s", name)
	}
	return ips[0], err
}

// GetMyIPSlice 获取指定网卡上的 IP 列表
func GetMyIPSlice(name string) ([]string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, errors.Wrapf(err, "GetMyIPSlice(%s) failed", name)
	}

	addrs, err := iface.Addrs()

	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			result = append(result, ipNet.IP.String())
		}
	}
	return result, nil
}
