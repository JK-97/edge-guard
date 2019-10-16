package cmd

const (
    InitPath             = "/edge/init"
  
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
    "docker-support-logging-go-arm64:1.0.1",
    "edgexfoundry/docker-support-logging-go-arm64:1.0.1",
    "edgexfoundry/docker-core-metadata-go-arm64:1.0.1",
    "edgexfoundry/docker-core-data-go-arm64:1.0.1",
    "edgexfoundry/docker-core-command-go-arm64:1.0.1",
}