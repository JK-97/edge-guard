package subprocess

import (
	"context"
	"github.com/JK-97/edge-guard/management/programmanage"

	log "github.com/JK-97/edge-guard/lowapi/logger"
)

type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
	Daemon        bool   `short:"d" long:"daemon" description:"run as daemon"`
	EnvFile       string `long:"env-file" description:"the environment file"`
}

// RunServer 启动 Edgenode 组件子进程
func RunServer(ctx context.Context) error {
	return runProcess(ctx, programmanage.GetJxConfig())
}

func runProcess(ctx context.Context, config string) error {
	// infinite loop for handling Restart ('reload' command)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		s := NewSupervisor(config)

		if sErr, addedGroup, changedGroup, removedGroup := s.Reload(); sErr != nil {
			return sErr
		} else {
			log.Info("addedGroup: ", addedGroup)
			log.Info("changedGroup: ", changedGroup)
			log.Info("removedGroup: ", removedGroup)
		}
		s.WaitForExit(ctx)
	}
}
