package register

import (
	"encoding/base64"
	"time"
)

const (
	// FallBackAuthHost 默认集群地址
	FallBackAuthHost = "http://auth.iotedge.jiangxingai.com:1054"

	wireguardRegisterPath = "/api/v1/wg/register"
	openvpnRegisterPath   = "/api/v1/openvpn/register"

	prefix = 512
	suffix = 128

	registerTimeout = time.Second * 10
)

var enc = base64.NewEncoding("ABCDEFGHIJKLMNOabcdefghijklmnopqrstuvwxyzPQRSTUVWXYZ0123456789-_").WithPadding(base64.NoPadding)
