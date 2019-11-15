package device

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	log "jxcore/go-utils/logger"
	"jxcore/lowapi/utils"
	"jxcore/version"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func GetDeviceType() (devicetype string) {
	return version.Type
}

func GetDevice() (device *Device, err error) {
	readdata, err := ioutil.ReadFile("/edge/init")
	utils.CheckErr(err)
	err = yaml.Unmarshal(readdata, &device)
	utils.CheckErr(err)
	return
}

// BuildWokerID 生成wokerid
func BuildWokerID() string {
	perfilx := "J"
	if runtime.GOARCH == "amd64" {
		perfilx = perfilx + "02"
	} else {
		perfilx = perfilx + "01"
	}
	var cpuInfoFile string = "/proc/cpuinfo"
	var GpsInfoScript string = "python /jxbootstrap/worker/scripts/G8100_NoMCU.py CMD AT+CGSN"
	var md5info [16]byte
	content, err := ioutil.ReadFile(cpuInfoFile)
	utils.CheckErr(err)
	if strings.Contains(string(content), "Serial") {
		md5info = md5.Sum(content[len(string(content))-17:])
	} else {
		for index := 0; index < 10; index++ {
			//小概率会获得空的数据,需重试
			gpsInfo, err := exec.Command("/bin/sh", "-c", GpsInfoScript).Output()
			utils.CheckErr(err)
			result := strings.ReplaceAll(string(gpsInfo), "\n", "")
			result = strings.TrimSpace(result)
			if len(result) > 10 {
				md5info = md5.Sum([]byte(result))
				break
			}

		}
	}

	md5str := fmt.Sprintf("%x", md5info)
	if len(md5str) < 7 {
		panic("can't generate workerid'")
	}
	workerid := perfilx + md5str[len(md5str)-7:]
	return workerid
}

// BuildDeviceInfo
func (d *Device) BuildDeviceInfo(vpnmodel Vpn, ticket string, authHost string) {
	if d == nil {
		d = new(Device)
	}
	if d.WorkerID == "" {
		d.WorkerID = BuildWokerID()
	}
	if vpnmodel == VPNModeRandom {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		vpnmodel = vpnSlice[r.Intn(len(vpnSlice))]
	}

	if GetDeviceType() == version.Pro {
		//pro
		//有dhcpserver则不再变动
		if d.DhcpServer != "" {
			d.DhcpServer = authHost
		} else {
			switch vpnmodel {
			case VPNModeLocal:
				d.DhcpServer = VPNModeLocal.String()
			case VPNModeWG, VPNModeOPENVPN, VPNModeRandom:
				d.DhcpServer = authHost
			default:
				log.Fatal("err vpnmodel")
			}
		}

		reqinfo := buildkeyreq{Workerid: d.WorkerID, Ticket: ticket}
		data, err := json.Marshal(reqinfo)
		if err != nil {
			log.Error(err)
		}
		//通过域名获取key
		body := bytes.NewBuffer(data)
		log.Info("Post to ", d.DhcpServer+BOOTSTRAPATH)

		// http.DefaultClient.Timeout = 8 * time.Second
		resp, err := http.Post(d.DhcpServer+BOOTSTRAPATH, "application/json", body)
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
		d.Key = respinfo.Data.Key
		d.Vpn = vpnmodel
		log.Info("Completed")
	} else {
		//base
		if vpnmodel != VPNModeLocal || authHost != VPNModeLocal.String() {
			log.Fatal("Base version can not support networking")
		}
		d.Vpn = VPNModeLocal
		d.DhcpServer = VPNModeLocal.String()

	}
	log.Info("Update Init Config File")
	outputdata, err := yaml.Marshal(d)
	utils.CheckErr(err)
	ioutil.WriteFile("/edge/init", outputdata, 0666)

}
