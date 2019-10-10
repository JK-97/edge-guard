package route

import (
    "net/http"
    "jxcore/web/controller"
)

// Routes adds routes to http
func Routes() http.Handler {
    mux := http.NewServeMux()
    handler := logRequest(mux)

    mux.HandleFunc("/api/v1/ping", controller.UpdateByDeb)

    return handler
}
