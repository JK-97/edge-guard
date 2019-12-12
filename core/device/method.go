package device

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/utils"
	"jxcore/version"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	initPath        string = "/edge/init"
	cpuInfoFile     string = "/proc/cpuinfo"
	GpsInfoScript   string = "python /jxbootstrap/worker/scripts/G8100_NoMCU.py CMD AT+CGSN"
	X86IdInfoScript string = "dmidecode | grep 'Serial Number' | head -1 | awk -F\":\" '{gsub(\" ^ \", \"\", $2); print $2}'"
)

var device *Device

func GetDeviceType() (devicetype string) {
	return version.Type
}

func getDevice() error {
	if _, err := os.Stat(initPath); err != nil {
		_ = ioutil.WriteFile(initPath, []byte{}, 0755)
	}
	readdata, _ := ioutil.ReadFile(initPath)
	err := yaml.Unmarshal(readdata, &device)
	if device == nil {
		device = &Device{}
	}
	return err
}

// GetDevice 获取节点信息，从/edge/init读取一次后存入cache
func GetDevice() (*Device, error) {
	var err error
	if device == nil {
		err = getDevice()
	}
	return device, err
}

// BuildWokerID 生成wokerid
func BuildWokerID() (string, error) {
	perfilx := "J"
	if runtime.GOARCH == "amd64" {
		perfilx = perfilx + "02"
	} else {
		perfilx = perfilx + "01"
	}

	var md5info = [16]byte{}
	switch runtime.GOARCH {
	case "amd64":
		err := buildX86ID(md5info[:])
		if err == nil {
			return "", err
		}
	case "arm64":
		err := buildRK3399ID(md5info[:])
		if err != nil {
			log.Error(err)
		} else {
			goto splice
		}
		err = buildNanoID(md5info[:])
		if err != nil {
			return "", err
		}
	}
splice:
	md5str := fmt.Sprintf("%x", md5info)
	if md5str == "0000000" {
		return "", errors.New("WorkerID Error")
	}
	workerid := perfilx + md5str[len(md5str)-7:]
	return workerid, nil
}

//SetHostname设置hostname
func (d *Device) SetHostname(workerid string) error {
	err := exec.Command("hostnamectl", "set-hostname", workerid).Run()
	return err
}

// BuildDeviceInfo 生成设备信息
func (d *Device) BuildDeviceInfo(vpnmodel Vpn, ticket string, authHost string) error {
	if d == nil {
		d = new(Device)
	}

	if vpnmodel == VPNModeRandom {
		//随机模式
		r := rand.New(rand.NewSource(time.Now().Unix()))
		vpnmodel = vpnSlice[r.Intn(len(vpnSlice))]
	}

	d.DhcpServer = authHost
	d.Vpn = vpnmodel

	respone, err := d.RequestKey(ticket)
	if err != nil {
		return err
	}

	if respone.Data.Key == "" {
		return errors.New("Request 获取 key 为空 ")
	}
	d.Key = respone.Data.Key

	err = d.UpdateDeviceInfo()
	if err != nil {
		return err
	}

	return nil
}

//RequestKey 通过ticket请求key
func (d *Device) RequestKey(ticket string) (*buildkeyresp, error) {
	request := buildkeyreq{Workerid: d.WorkerID, Ticket: ticket}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//通过域名服务器获取key
	body := bytes.NewBuffer(data)

	resp, err := http.Post(d.DhcpServer+BOOTSTRAPATH, "application/json", body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := buildkeyresp{}
	err = json.Unmarshal(respdata, &response)
	return &response, err
}

// UpdateDeviceInfo 更新device配置
func (d *Device) UpdateDeviceInfo() error {

	outputdata, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/edge/init", outputdata, 0666)
	if err != nil {
		return err
	}
	return nil
}

func buildX86ID(md5info []byte) error {
	if runtime.GOARCH == "amd64" {
		data, err := exec.Command("/bin/bash", "-c", X86IdInfoScript).Output()
		if err != nil {
			return err
		}

		info := md5.Sum(data)
		copy(md5info, info[:])
	}
	return errors.New("This Platform is not amd64")
}

func buildRK3399ID(md5info []byte) error {
	content, err := ioutil.ReadFile(cpuInfoFile)
	if err != nil {
		return err
	}
	if !strings.Contains(string(content), "Serial") {
		// RK品台
		return errors.New("This Platform is not rk")
	}
	info := md5.Sum(content[len(string(content))-17:])
	copy(md5info, info[:])
	return nil
}

func buildNanoID(md5info []byte) error {
	for index := 0; index < 10; index++ {
		//小概率会获得空的数据,需重试
		gpsInfo, err := exec.Command("/bin/sh", "-c", GpsInfoScript).Output()
		utils.CheckErr(err)
		result := strings.ReplaceAll(string(gpsInfo), "\n", "")
		result = strings.TrimSpace(result)
		if len(result) > 10 {
			info := md5.Sum([]byte(result))
			copy(md5info, info[:])
		}
	}
	return nil
}
