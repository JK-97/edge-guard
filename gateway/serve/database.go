package serve

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/JK-97/edge-guard/gateway/option"
)

type DatabaseHandler struct {
	option.DatabaseConfig
	http.Handler
}

// NewDatabaseHandler 获取新的 DatabaseHandler 实例
func NewDatabaseHandler(config option.DatabaseConfig) http.Handler {
	rpURL, err := url.Parse(config.Host)
	if err != nil {
		panic(err)
	}
	rpURL.Path = "/api/v1/" + rpURL.Path
	return &DatabaseHandler{
		DatabaseConfig: config,
		Handler:        httputil.NewSingleHostReverseProxy(rpURL),
	}
}
