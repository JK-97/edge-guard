package web

import (
	"context"

	"net/http"
	"time"

	"jxcore/lowapi/logger"
)

func Serve(ctx context.Context, addr string, handler http.Handler, wait time.Duration) error {
	srv := http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	return srv.Shutdown(ctx)
}
