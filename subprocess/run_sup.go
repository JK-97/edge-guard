package subprocess

import (
	"context"
	"jxcore/management/programmanage"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
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

// RunMcuProcess 启动 mcu 相关子进程
func RunMcuProcess(ctx context.Context) error {
	return runProcess(ctx, programmanage.GetMcuConfig())
}

// RunJxserving 启动 Jxserving
func RunJxserving(ctx context.Context) error {
	return runProcess(ctx, programmanage.GetJxserving())
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
