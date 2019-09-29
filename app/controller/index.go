package controller

import (
	"fmt"
	"net/http"
)

//返回pong
func PingGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// 获取html主叶
func IndexGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveStatic("web/template/index.html", w, r)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

