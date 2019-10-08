package route

import (
    "jxcore/web/controller/index"
    "net/http"
)

// Routes adds routes to http
func Routes() http.Handler {
    mux := http.NewServeMux()
    handler := logRequest(mux)

    mux.HandleFunc("/api/v1/ping", index.PingGet)

    return handler
}
