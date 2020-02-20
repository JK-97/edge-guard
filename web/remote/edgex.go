package remote

import (
	"context"
	"fmt"
	"io/ioutil"
	"jxcore/lowapi/logger"
	"jxcore/web/controller/utils"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	urlGetDeviceServices = "http://edgegw.iotedge:48081/api/v1/deviceservice"
	edgexNetworkName     = "edgex"
)

var edgexOfficalContainerNames = map[string]bool{
	"edgex-mongo":           true,
	"edgex-files":           true,
	"edgex-core-consul":     true,
	"edgex-core-command":    true,
	"edgex-core-metadata":   true,
	"edgex-core-data":       true,
	"edgex-support-logging": true,
}

type DeviceService struct {
	Name        string
	Description string
	ProxyName   string // 用于路由请求到 device service
	IP          string
	Port        int
	Version     string
	Alive       bool
}

var cacheProxyNameToDS = map[string]*DeviceService{}

// 获取edgex device service 信息，步骤:
// 1. 获取edgex device service name, hostname, port
// 2. 通过docker api 获取edgex 网络下所有 device service 对应关系 hostname -> ip
// 3. 访问 ip:port/api/v1/ping 获取版本号
// 4. 访问 ip:port/api/v1/name 获取路由名称
func GetDeviceServices(ctx context.Context) ([]*DeviceService, error) {
	edgexDSs, err := getEdgexDeviceServices()
	if err != nil {
		return nil, err
	}

	edgexIPs, err := getEdgexIPs(ctx)
	if err != nil {
		return nil, err
	}

	deviceServices := []*DeviceService{}
	for _, eds := range edgexDSs {
		ds := &DeviceService{Name: eds.Name, Description: eds.Description}
		deviceServices = append(deviceServices, ds)

		if ip, ok := edgexIPs[eds.Addressable.Address]; ok {
			ds.IP = ip
			ds.Port = eds.Addressable.Port

			ds.updateDeviceServicePing()
			ds.updateDeviceServiceProxyName()

			if ds.ProxyName != "" {
				cacheProxyNameToDS[ds.ProxyName] = ds
			}
		}
	}
	return deviceServices, nil
}

var ErrNotFound = fmt.Errorf("not found")

func GetDeviceServiceByProxyName(ctx context.Context, proxyName string) (*DeviceService, error) {
	if ds, ok := cacheProxyNameToDS[proxyName]; ok {
		return ds, nil
	}
	_, err := GetDeviceServices(ctx)
	if err != nil {
		return nil, err
	}
	if ds, ok := cacheProxyNameToDS[proxyName]; ok {
		return ds, nil
	}
	return nil, ErrNotFound
}

const edgexVersionNotReplaced = "to be replaced by makefile"

// updateDeviceServicePing 更新 device service 连接状态 和 版本号
func (ds *DeviceService) updateDeviceServicePing() {
	url := fmt.Sprintf("http://%v:%v/api/v1/ping", ds.IP, ds.Port)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	alive := true

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	version := string(data)

	if version == edgexVersionNotReplaced {
		version = ""
	}

	ds.Alive = alive
	ds.Version = version
}

// updateDeviceServiceProxyName 更新 device service 的路由名称
func (ds *DeviceService) updateDeviceServiceProxyName() {
	url := fmt.Sprintf("http://%v:%v/api/v1/name", ds.IP, ds.Port)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && resp.StatusCode == http.StatusOK {
		ds.ProxyName = string(data)
	}
}

type edgexDeviceServiceResp struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Addressable struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"addressable"`
}

func getEdgexDeviceServices() ([]edgexDeviceServiceResp, error) {
	resp, err := http.Get(urlGetDeviceServices)
	if err != nil {
		return nil, err
	}

	var deviceServices []edgexDeviceServiceResp
	err = utils.UnmarshalJson(resp.Body, &deviceServices)
	return deviceServices, err
}

// getEdgexIPs gets edgex device services ip, returns <ds hostname> -> <ip>
func getEdgexIPs(ctx context.Context) (map[string]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	network, err := cli.NetworkInspect(ctx, edgexNetworkName, types.NetworkInspectOptions{})
	if err != nil {
		return nil, err
	}

	edgexIPs := map[string]string{}
	for _, c := range network.Containers {
		if _, ok := edgexOfficalContainerNames[c.Name]; !ok {
			container, err := cli.ContainerInspect(ctx, c.Name)
			if err != nil {
				logger.Error(err)
				continue
			}

			networkConfig, ok := container.NetworkSettings.Networks[edgexNetworkName]
			if !ok {
				logger.Errorf("Container %s not in edgex network, while network inspect shows it's in.", c.Name)
				continue
			}
			for _, a := range networkConfig.Aliases {
				edgexIPs[a] = networkConfig.IPAddress
			}
		}
	}
	return edgexIPs, nil
}
