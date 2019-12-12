package route

import (
	"errors"
	"jxcore/lowapi/logger"
	"jxcore/web/controller"
	"net/http"
)

func logRequestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case error:
					err = x
				case string:
					err = errors.New(x)
				default:
					err = errors.New("Unknown panic")
				}
				controller.Error(w, err, 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
