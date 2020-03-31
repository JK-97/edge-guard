package serve

import (
	// "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/JK-97/edge-guard/gateway/log"
	"github.com/JK-97/edge-guard/gateway/option"
)

type service struct {
	option.Service
	proxyURL *url.URL
}

type route struct {
	option.Route
	matcher *regexp.Regexp
}

// ReverseProxyHandler 反向代理服务
type ReverseProxyHandler struct {
	option.ProxyServerConfig
	// ReverseProxy string
	Proxy    *url.URL
	Services *map[string]service
	Routes   *[]route

	masterHandler http.Handler
}

// CreateReverseProxyFromOpion 根据配置文件中生成反向代理服务
func CreateReverseProxyFromOpion(opt *option.ProxyServerConfig) *ReverseProxyHandler {

	result, err := url.Parse(opt.MasterProxy)
	if err != nil {
		log.Println(err)
	}

	services := make(map[string]service, len(opt.Services))
	for key, srv := range opt.Services {
		srv2 := service{Service: srv}
		if !srv2.MasterOnly {
			result, err := url.Parse(srv.Proxy)
			if err != nil {
				log.Println(err)
			}
			srv2.proxyURL = result
		}
		services[key] = srv2
	}

	routes := make([]route, len(opt.Routes))
	for i, r := range opt.Routes {
		result, err := regexp.Compile(r.Matcher)

		if err != nil {
			log.Println(err)
		}
		routes[i] = route{
			Route:   r,
			matcher: result,
		}
	}

	handler := ReverseProxyHandler{
		ProxyServerConfig: *opt,
		Services:          &services,
		Routes:            &routes,
		Proxy:             result,

		masterHandler: httputil.NewSingleHostReverseProxy(result),
	}

	return &handler
}

// ServeMasterHTTP 把请求全部转发到 Master 上
func (h *ReverseProxyHandler) ServeMasterHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := h.masterHandler
	r.Host = h.Proxy.Host
	log.Printf("Proxy: %s\n", h.Proxy)

	proxy.ServeHTTP(w, r)
}

func (h *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localOnly := false

	if r.Header.Get("X-To-Master") != "" {
		h.ServeMasterHTTP(w, r)
		return
	} else if r.Header.Get("X-Route-Gateway") != "" {
		localOnly = true
	}

	for _, route := range *h.Routes {
		if route.matcher.MatchString(r.URL.Path) {
			service, err := (*h.Services)[route.Name]
			if !err {
				continue
			}
			var _url *url.URL
			if localOnly {
				_url = service.proxyURL
			} else if service.MasterOnly {
				h.ServeMasterHTTP(w, r)
				return
			} else {
				_url = service.proxyURL
			}
			proxy := httputil.NewSingleHostReverseProxy(_url)
			// proxy.Transport = &proxyRoundTrip{Transport: *http.DefaultTransport.(*http.Transport)}
			r.URL.Path = strings.Replace(r.URL.Path, "/internalapi", "/api", 1)

			log.Printf("Proxy: %s\n", _url)
			proxy.ServeHTTP(w, r)
			return
		}
	}
	log.Printf("Not Found.\n")
	Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
