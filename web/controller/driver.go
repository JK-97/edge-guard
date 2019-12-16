package controller

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/lowapi/logger"
	"net/http"
	"regexp"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

const (
	urlGetDeviceServices = "http://edgegw.iotedge:48081/api/v1/deviceservice"
)

type edgexDeviceService struct {
	Name string `json:"name"`
}

type GetEdgexDriversRespDriver struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Alive       bool   `json:"alive"`
}

type GetEdgexDriversResp struct {
	Driver []GetEdgexDriversRespDriver `json:"drivers"`
}

// 获取edgex device service 信息
// 步骤: 1. 获取edgex device service name
// 2. 获取health，解析ping 接口返回的版本号
func GetEdgexDrivers(w http.ResponseWriter, r *http.Request) {
	cli, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	deviceServices, err := getDeviceServices()
	if err != nil {
		panic(err)
	}

	resp := GetEdgexDriversResp{}
	for _, ds := range deviceServices {
		healthChecks, _, err := cli.Health().Checks(ds.Name, nil)
		if err != nil {
			logger.Errorf("Failed to get health of %v: %v", ds.Name, err)
			continue
		}
		driver := GetEdgexDriversRespDriver{Name: ds.Name}
		for _, hc := range healthChecks {
			if hc.Status == "passing" {
				driver.Alive = true
				// version 在最后一个单词
				version := hc.Output[strings.LastIndex(hc.Output, " ")+1:]
				if versionRegex.Match([]byte(version)) {
					driver.Version = version
				}
			}
		}
		resp.Driver = append(resp.Driver, driver)
	}
	RespondSuccessJSON(resp, w)
}

func getDeviceServices() ([]edgexDeviceService, error) {
	resp, err := http.Get(urlGetDeviceServices)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceServices []edgexDeviceService
	err = json.Unmarshal(data, &deviceServices)
	return deviceServices, err
}

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile("[0-9]+(.[0-9]+)*")
}
