package device

// DeviceKey 设备信息
var DeviceInstance = Device{}

const (
    // FallBackAuthHost 默认集群地址
    FallBackAuthHost      string = "http://auth.iotedge.jiangxingai.com:1054"
    bootstraPath          string = "/api/v1/bootstrap"
    wireguardRegisterPath string = "/api/v1/wg/register"
    openvpnRegisterPath   string = "/api/v1/openvpn/register"
)

// VPN 模式
const (
    VPNModeRandom  string = "random"
    VPNModeWG      string = "wireguard"
    VPNModeOPENVPN string = "openvpn"
    VPNModeLocal   string = "local"
)

var vpnSlice []string = []string{VPNModeWG, VPNModeOPENVPN}
