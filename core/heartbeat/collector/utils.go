package collector

import (
	"github.com/JK-97/edge-guard/lowapi/logger"
	"strconv"
	"strings"
)

// hexToInt 字符串形式十六进制数转 int
func hexToInt(hex string) (int, error) {
	var r int

	for _, it := range hex[2:] {
		r <<= 4
		r += int(it - '0')
	}

	return r, nil
}

func mustInt(s string) int {
	s = strings.TrimSpace(s)
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn(err, "Input:", s)
	}
	return i
}

func translateFromHuman(value string) int64 {
	segments := strings.Split(value, " ")

	var r int
	r, _ = strconv.Atoi(segments[0])
	rr := int64(r)
	if len(segments) == 1 {
		return rr
	}
	switch segments[1] {
	case "kB":
		rr = rr << 10
	default:
		logger.Warn(value)
	}

	return rr
}

var virtualNetworkPrefixes = []string{
	"docker",
	"veth",
	"br-",
}

// 判断是否为虚拟网卡或网桥
func isVirtualNetwork(name string) bool {
	for _, it := range virtualNetworkPrefixes {

		if strings.HasPrefix(name, it) {
			return true
		}
	}
	return false
}
