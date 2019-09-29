package monitor
//
//import (
//	"jxcore/log"
//
//	"github.com/google/gousb"
//)
////usb的监控
//func UsbMonitor() {
//
//
//	ctx := gousb.NewContext()
//	defer ctx.Close()
//
//	// 返回所有连接设备的pid与id
//	dev, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
//		return true
//	})
//	if err != nil {
//		log.Fatalf("Could not open a device: %v", err)
//	}
//	for _, perdev := range dev {
//		log.Info(perdev)
//		perdev.Close()
//	}
//
//}
