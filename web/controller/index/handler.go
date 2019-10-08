package index

import (
    "fmt"
    "net/http"
)

//返回pong
func PingGet(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "pong")
}

