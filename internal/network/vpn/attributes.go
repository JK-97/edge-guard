package vpn

import (
	"errors"
	"time"
)

const (
	WireGuardInterface = "wg0"
	OpenVPNInterface   = "tun0"

	wgWaitTimeout      time.Duration = 20 * time.Second
	openvpnWaitTimeout time.Duration = 20 * time.Second

	openvpnSuccessMessage = "Initialization Sequence Completed"
	openvpnConfigPath     = "/etc/openvpn/client.ovpn"

	wireguardConfigDir = "/etc/wireguard/"
	openvpnConfigDir   = "/etc/openvpn/"
)

var (
	// Open VPN 启动超时
	ErrOpenVPNTimeout = errors.New("Start Open VPN Timeout")

	myVpnIP string
)
