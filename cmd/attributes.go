package cmd

const (
	LogsPath           = "/edge/logs"
	TargetVersionFile  = "/etc/edgetarget"
	CurrentVersionFile = "/edge/VERSION"
)

var DependsImages = []string{
	"registry.jiangxingai.com:5000/cadvisor:arm64v8-cpu-0.1.0",
	"registry.jiangxingai.com:5000/edgex/device-service/camera:arm64v8-cpu-1.1.1",
	"registry.jiangxingai.com:5000/nginx-rtmp:arm64v8-cpu-0.1.0",
	"registry.jiangxingai.com:5000/device-statsite:arm64v8_cpu_1.0.0",
	"registry.jiangxingai.com:5000/config-agent:arm64v8_cpu_1.0.3",
	"registry.jiangxingai.com:5000/tensorflow-serving:arm64v8_cpu_0.1.0",
	"registry.jiangxingai.com:5000/consul:arm64v8_cpu_1.5.3",
	"edgexfoundry/docker-edgex-volume-arm64:1.0.0",
	"edgexfoundry/docker-core-config-seed-go-arm64:1.0.0",
	"edgexfoundry/docker-support-logging-go-arm64:1.0.1",
	"edgexfoundry/docker-core-metadata-go-arm64:1.0.1",
	"edgexfoundry/docker-core-data-go-arm64:1.0.1",
	"edgexfoundry/docker-core-command-go-arm64:1.0.1",
}

var BinFilesMAP = map[string]string{
	"watchdog":        "/edge/tools/mcutools/watchdog/bin/watchdog",
	"powermanagement": "/edge/tools/mcutools/powermanagement/bin/powermanagement",
	"db":              "/edge/mnt/db/bin/db",
	"mcuserver":       "/edge/tools/mcutools/mcuserver/bin/mcuserver",
	"telegraf":        "/edge/monitor/telegraf/bin/telegraf",
	"jxserving":       "/jxserving/run.py",
	"filelistener":    "/edge/tools/nodetools/filelistener/bin/filelistener",
	"cleaner":         "/edge/tools/nodetools/cleaner/bin/cleaner",
}

var DependsFile = []string{
	"docker",
	"supervisord",
	"openvpn",
	"docker-compose",
	"aptitude",
	"dnsmasq",
}
