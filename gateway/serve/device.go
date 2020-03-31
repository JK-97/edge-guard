package serve

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/JK-97/edge-guard/gateway/option"
)

type DeviceHandler struct {
	option.DeviceConfig
	http.Handler
}

// NewDeviceHandler 获取新的 DeviceHandler 实例
func NewDeviceHandler(config option.DeviceConfig) http.Handler {
	director := func(req *http.Request) {
		edgexModuleName := strings.SplitN(req.URL.Path, "/", 2)[0]
		host := config.Hosts[edgexModuleName]

		u, err := url.Parse(host)
		if err != nil {
			panic(err)
		}

		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = strings.TrimPrefix(req.URL.Path, edgexModuleName)
	}

	return &DeviceHandler{
		DeviceConfig: config,
		Handler:      &httputil.ReverseProxy{Director: director},
	}
}
