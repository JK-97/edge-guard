package register

const (
    // FallBackAuthHost 默认集群地址
    FallBackAuthHost      string = "http://auth.iotedge.jiangxingai.com:1054"
    bootstraPath          string = "/api/v1/bootstrap"
    wireguardRegisterPath string = "/api/v1/wg/register"
    openvpnRegisterPath   string = "/api/v1/openvpn/register"
)
