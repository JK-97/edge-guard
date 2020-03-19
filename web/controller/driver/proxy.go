package driver

import (
	"context"
	"fmt"
	"jxcore/web/remote"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dsname := vars["dsname"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	deviceService, err := remote.GetDeviceServiceByProxyName(ctx, dsname)
	if err != nil {
		panic(err)
	}

	target, err := url.Parse(fmt.Sprintf("http://%v:%v", deviceService.IP, deviceService.Port))
	if err != nil {
		panic(err)
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.ServeHTTP(w, r)
}
