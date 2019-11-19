package web

import (
	"context"
	"jxcore/gateway/log"
	"net/http"
	"time"
)

func Serve(ctx context.Context, port string, handler http.Handler, wait time.Duration) error {
	srv := http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	return srv.Shutdown(ctx)
}
