package device

// DeviceKey 设备信息
var DeviceInstance = Device{}


// VPN 模式
const (
    VPNModeRandom  string = "random"
    VPNModeWG      string = "wireguard"
    VPNModeOPENVPN string = "openvpn"
    VPNModeLocal   string = "local"
)

var vpnSlice []string = []string{VPNModeWG, VPNModeOPENVPN}
