package regeister

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"jxcore/component"
	"jxcore/log"
	"jxcore/utils"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	WorkID     string `json:"workerid"`
	Key        string `json:"key"`
	DhcpServer string `json:"dhcpserver"`
	Vpn        string `json:"vpn"`
}

type buildkeyreq struct {
	Workerid string `json:"wid"`
	Ticket   string `json:"ticket"`
}

type buildkeyresp struct {
	Data data `json:"data"`
}
type data struct {
	Key         string `json:"key"`
	DeadLine    string `json:"deadLine"`
	RemainCount string `json:"remainCount"`
}

// DeviceKey 设备信息
var DeviceKey = DeviceInfo{}

const (
	// FallBackAuthHost 默认集群地址
	FallBackAuthHost      string = "http://auth.iotedge.jiangxingai.com:1054"
	bootstraPath          string = "/api/v1/bootstrap"
	wireguardRegisterPath string = "/api/v1/wg/register"
	openvpnRegisterPath   string = "/api/v1/openvpn/register"
)

// var bootstrapurl = "http://auth.iotedge.jiangxingai.com:1054/api/v1/bootstrap"
// var wireguardregisterurl = "http://auth.iotedge.jiangxingai.com:1054/api/v1/wg/register"
// var openvpnregisterurl = "http://auth.iotedge.jiangxingai.com:1054/api/v1/openvpn/register"

var Connectable = make(chan bool, 1)

// Patternmatching 模式匹配 校验
func Patternmatching() {
	if !utils.Exists("/edge/init") {

		log.Warn("The current host has not been initialized")
	}
	log.WithFields(log.Fields{"Operating": "VpnPatternMatching"}).Info("Device initialized")
	deviceinfo, err := ioutil.ReadFile("/edge/init")
	if err != nil {
		log.Error(err)
	}
	yaml.Unmarshal(deviceinfo, &DeviceKey)
	if utils.Exists("/edge/mnt/mq") {

	} else {
		if DeviceKey.Vpn != VPNModeLocal {
			panic("Vpn mode does not match the file")
		}
	}
	switch DeviceKey.Vpn {
	case VPNModeLocal:
		log.WithFields(log.Fields{"Operating": "VpnPatternMatching"}).Info("Local mode")
		Connectable <- true
		// not thing to do
		time.Sleep(3 * time.Second)
		component.StopComponent("vpn")
	case VPNModeWG, VPNModeOPENVPN:
		if DeviceKey.DhcpServer == "" || DeviceKey.Key == "" || DeviceKey.WorkID == "" {
			log.Fatal("Missing current device information. Please run jxcore bootstrap to generate the device before running.")
		} else {
			//设备信息正常
			r := regexp.MustCompile(".*")
			regexpres := r.FindString(DeviceKey.DhcpServer)
			if regexpres != "" {
				log.WithFields(log.Fields{"Operating": "VpnPatternMatching"}).Info("using the ", DeviceKey.Vpn, " model")
				// 注册
				register(DeviceKey.WorkID, DeviceKey.Key)
			} else {
				log.WithFields(log.Fields{"Operating": "VpnPatternMatching"}).Info("Registration server format error")
			}
		}
	default:
		log.Error("Error VPN model")
	}

}

// ReadDeviceInfo 读取设备信息
func ReadDeviceInfo() (info DeviceInfo, err error) {
	readdata, err := ioutil.ReadFile("/edge/init")
	if err != nil {
		log.Error(err)
		return
	}
	err = yaml.Unmarshal(readdata, &info)
	return
}

// VPN 模式
const (
	VPNModeRandom  string = "random"
	VPNModeWG      string = "wireguard"
	VPNModeOPENVPN string = "openvpn"
	VPNModeLocal   string = "local"
)

var vpnSlice []string = []string{VPNModeWG, VPNModeOPENVPN}

// // DailVPN 连接 VPN
// func DailVPN(workId, vpnmodel, authHost, key string) error {

// 	return nil
// }

// BuildDeviceInfo 生成key
func BuildDeviceInfo(vpnmodel string, ticket string, authHost string) {
	deviceinfo, err := ReadDeviceInfo()
	if err != nil {
		log.Error(err)
	}
	if authHost == "" {
		authHost = FallBackAuthHost
	}

	//有workerid 则不再变动,
	if deviceinfo.WorkID != "" {

	} else {
		deviceinfo.WorkID = BuildWokerID()
	}
	if vpnmodel == VPNModeRandom {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		vpnmodel = vpnSlice[r.Intn(len(vpnSlice))]
	}
	//有dhcpserver则不再变动
	if deviceinfo.DhcpServer != "" {
		deviceinfo.DhcpServer = authHost
	} else {
		switch vpnmodel {
		case VPNModeLocal:
			deviceinfo.DhcpServer = "local"
		case VPNModeWG, VPNModeOPENVPN, VPNModeRandom:
			deviceinfo.DhcpServer = authHost
		default:
			log.Fatal("err vpnmodel")
		}
	}

	if deviceinfo.DhcpServer != "local" {
		//联网模式

		reqinfo := buildkeyreq{Workerid: deviceinfo.WorkID, Ticket: ticket}
		data, err := json.Marshal(reqinfo)
		if err != nil {
			log.Error(err)
		}
		//通过域名获取key
		body := bytes.NewBuffer(data)
		log.Info("Post to ", deviceinfo.DhcpServer+bootstraPath)

		// http.DefaultClient.Timeout = 8 * time.Second
		resp, err := http.Post(deviceinfo.DhcpServer+bootstraPath, "application/json", body)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Status:", resp.Status)
		respdata, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		respinfo := buildkeyresp{}
		json.Unmarshal(respdata, &respinfo)
		deviceinfo.Key = respinfo.Data.Key
		deviceinfo.Vpn = vpnmodel
		log.Info("Completed")
	} else {
		//本地模式
		deviceinfo.Key = "local"
		deviceinfo.Vpn = vpnmodel
	}
	log.Info("Update Init Config File")
	outdata, err := yaml.Marshal(deviceinfo)
	if err != nil {
		log.Error(err)
	}
	ioutil.WriteFile("/edge/init", outdata, 0666)
}
