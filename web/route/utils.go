package route
import (
    log "jxcore/go-utils/logger"
    "net/http"
)

func logRequest(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)
    })
}

