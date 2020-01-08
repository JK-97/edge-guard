package utils

import (
	"context"
	"jxcore/lowapi/logger"
	"reflect"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

func CheckErr(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func RunAndLogError(f func() error) {
	err := f()
	if err != nil {
		logger.Error(err)
	}
}

func RunUntilSuccess(ctx context.Context, retryInterval time.Duration, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return err
		case <-ticker.C:
			err = f()
			if err == nil {
				return nil
			}
		}
	}
}

const (
	funcRestartWait = time.Second
)

// GoAndRestartOnError runs function in background associated with error group, until context done. Always restart if failed.
func GoAndRestartOnError(ctx context.Context, errGroup *errgroup.Group, name string, f func() error) {
	errGroup.Go(func() error { return RunAndRestartOnError(ctx, name, f) })
}

// RunAndRestartOnError runs function until context done. Always restart if failed.
func RunAndRestartOnError(ctx context.Context, name string, f func() error) error {
	for {
		logger.Infof("starting %s", name)
		err := f()
		if err != nil {
			logger.Errorf("%s stopped: %v", name, err)
		}

		select {
		case <-ctx.Done():
			return err
		default:
		}

		logger.Infof("%s will restart in %v", name, funcRestartWait)
		time.Sleep(funcRestartWait)
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
