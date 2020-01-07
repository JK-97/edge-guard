package core

import (
	"context"
	"jxcore/config/yaml"
	"jxcore/core/register"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/iface"
	"jxcore/management/updatemanage"

	"golang.org/x/sync/errgroup"
)

func ConfigSupervisor() {
	startupProgram := yaml.Config
	yaml.ParseAndCheck(*startupProgram, "")
}

func MaintainNetwork(ctx context.Context, noUpdate bool) error {
	dns.TrySetupDnsConfig()

	errGroup := errgroup.Group{}

	// 按优先级切换网口
	errGroup.Go(func() error { return iface.MaintainBestIFace(ctx) })

	// 第一次连接master成功，检查固件更新
	onFirstConnect := func() {
		manager := updatemanage.NewUpdateManager()
		manager.ReportVersion()
		if !noUpdate {
			manager.Start()
		}
	}
	errGroup.Go(func() error { return register.MaintainMasterConnection(ctx, onFirstConnect) })

	return errGroup.Wait()
}
