package web

import (
	"context"

	"net/http"
	"time"
)

func Serve(ctx context.Context, addr string, handler http.Handler, wait time.Duration) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	srv := http.Server{
		Addr:        addr,
		Handler:     handler,
	}
	var err error
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			cancel()
		}
	}()
	<-cancelCtx.Done()

	if err != nil {
		return err
	}

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), wait)
	defer cancelShutdown()
	return srv.Shutdown(ctxShutdown)
}
