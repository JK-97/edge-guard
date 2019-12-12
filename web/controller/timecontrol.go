package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jxcore/internal/config"
	"jxcore/lowapi/store/filestore"
	"jxcore/lowapi/system"
	"net/http"
	"time"

	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

const timeConfigKey = "time_config"

type timeRequest struct {
	Time int64 `json:"time"`
}

func SetTime(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	request := timeRequest{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		panic(err)
	}

	//设置日期时间，写入BIOS，以免重启失效
	err = system.RunCommand(fmt.Sprintf("date -u --set=@%f && hwclock -w", float64(request.Time)/1e9))
	if err != nil {
		panic(err)
	}
	RespondJSON(nil, w, 200)
}

func GetTime(w http.ResponseWriter, r *http.Request) {
	resp := timeRequest{Time: time.Now().UnixNano()}
	RespondJSON(resp, w, 200)
}

type ntpConfig struct {
	Enabled    bool   `json:"enabled"`
	ServerAddr string `json:"server_addr"`
}
type timeConfig struct {
	TimeZONE string    `json:"timezone"`
	Ntp      ntpConfig `json:"ntp"`
}

var defaultNtpConfData []byte
var defaultNtpConf = timeConfig{
	TimeZONE: "Asia/Shanghai",
	Ntp: ntpConfig{
		Enabled:    true,
		ServerAddr: "0.arch.pool.ntp.org",
	},
}

func init() {
	defaultNtpConfData, _ = json.Marshal(defaultNtpConf)
}

// timedatectl 用于配置时区，timesyncd 同步进程配置。文档：
// https://www.freedesktop.org/software/systemd/man/timedatectl.html

func GetNtpConfig(w http.ResponseWriter, r *http.Request) {
	data, err := filestore.KV.GetDefault(timeConfigKey, defaultNtpConfData)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}
}

func SetNtpConfig(w http.ResponseWriter, r *http.Request) {
	patch, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	oldData, err := filestore.KV.GetDefault(timeConfigKey, defaultNtpConfData)
	if err != nil {
		panic(err)
	}
	old, new, newData, err := patchConfig(oldData, patch)
	if err != nil {
		RespondReasonJSON(err, w, "Config not valid", 400)
		return
	}
	err = filestore.KV.Set(timeConfigKey, newData)
	if err != nil {
		panic(err)
	}

	if old.TimeZONE != new.TimeZONE {
		err := system.RunCommand("timedatectl set-timezone " + new.TimeZONE)
		if err != nil {
			panic(err)
		}
	}

	config.TimdsyncdCfg(new.Ntp.ServerAddr)

	_ = system.RunCommand("timedatectl set-ntp false")
	if new.Ntp.Enabled {
		err := system.RunCommand("timedatectl set-ntp true")
		if err != nil {
			panic(err)
		}
	}

	RespondJSON(nil, w, 200)
}

func patchConfig(oldData, patch []byte) (old, new timeConfig, newData []byte, err error) {
	newData, err = jsonpatch.MergePatch(oldData, patch)
	if err != nil {
		return
	}

	err = json.Unmarshal(oldData, &old)
	if err != nil {
		return
	}

	err = json.Unmarshal(newData, &new)
	return
}
