package register

import "encoding/base64"

const (
	// FallBackAuthHost 默认集群地址
	FallBackAuthHost string = "http://auth.iotedge.jiangxingai.com:1054"

	wireguardRegisterPath string = "/api/v1/wg/register"
	openvpnRegisterPath   string = "/api/v1/openvpn/register"
)

var prefix = 512
var suffix = 128
var enc = base64.NewEncoding("ABCDEFGHIJKLMNOabcdefghijklmnopqrstuvwxyzPQRSTUVWXYZ0123456789-_").WithPadding(base64.NoPadding)
