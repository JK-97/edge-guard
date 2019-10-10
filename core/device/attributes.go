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
const BOOTSTRAPATH  string = "/api/v1/bootstrap"
var vpnSlice []string = []string{VPNModeWG, VPNModeOPENVPN}
