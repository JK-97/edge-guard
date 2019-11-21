package route

import (
	"jxcore/web/controller"
	"net/http"
)

// Routes adds routes to http
func Routes() http.Handler {
	mux := http.NewServeMux()
	handler := logRequest(mux)

	mux.HandleFunc("/edgenode/exec/update", controller.UpdateByDeb)

	return handler
}
