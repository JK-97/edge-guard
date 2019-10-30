package serve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"jxcore/gateway/power"
)

type powerOffParam struct {
	DelayTime int // 多久后关机，单位： 秒，至少为 5 秒
	WakeTime  int // 多久后开机，单位： 秒，至少为 60 秒
}

// PowerOffHTTP 关机
func PowerOffHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}

	p := powerOffParam{}
	buff, err := ioutil.ReadAll(r.Body)
	if err == nil {
		json.Unmarshal(buff, &p)
	}
	if p.DelayTime <= 0 && p.DelayTime != -1 {
		p.DelayTime = 5
	}
	if p.WakeTime <= 60 && p.WakeTime != -1 {
		p.WakeTime = 60
	}

	power.SystemPowerOff(p.DelayTime, p.WakeTime)
	data := map[string]interface{}{
		"DelayTime": p.DelayTime,
		"WakeTime":  p.WakeTime,
	}
	WriteData(w, &data)
}

// StartUpMode 开机模式
func StartUpMode(w http.ResponseWriter, r *http.Request) {

	mode := power.GetStartUpMode()
	data := map[string]interface{}{
		"Mode": int(mode),
	}
	WriteData(w, &data)
}
