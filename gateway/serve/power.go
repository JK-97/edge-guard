package serve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"jxcore/gateway/power"
)

type powerOffParam struct {
	WakeTime int
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
	if p.WakeTime <= 0 {
		p.WakeTime = 30
	}

	power.SystemPowerOff(p.WakeTime)
	data := map[string]interface{}{
		"WakeTime": p.WakeTime,
	}
	WriteData(w, &data)
}
