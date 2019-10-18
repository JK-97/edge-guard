package device

import "errors"

// DeviceKey 设备信息
var DeviceInstance = Device{}

// VPN 模式
const (
    VPNModeRandom  Vpn = "random"
    VPNModeWG      Vpn = "wireguard"
    VPNModeOPENVPN Vpn = "openvpn"
    VPNModeLocal   Vpn = "local"
)
const BOOTSTRAPATH string = "/api/v1/bootstrap"

var vpnSlice []Vpn = []Vpn{VPNModeWG, VPNModeOPENVPN}

func (v Vpn) String() string {
    switch v {
    case VPNModeWG:
        return "wireguard"
    case VPNModeOPENVPN:
        return "openvpn"
    case VPNModeRandom:
        return "random"
    case VPNModeLocal:
        return "local"
    default:
        return ""
    }
}

func (v Vpn) Interface() (string, error) {
    switch v {
    case VPNModeWG:
        return "wg0", nil
    case VPNModeOPENVPN:
        return "tun0", nil
    default:
        return "unknow", errors.New("no supported " + v.String())
    }

}
