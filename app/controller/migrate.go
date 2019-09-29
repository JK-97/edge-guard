package controller

import (
	"jxcore/utils"

	"io/ioutil"
	"jxcore/log"
	"jxcore/regeister"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Migrate 设备迁移
func Migrate(w http.ResponseWriter, r *http.Request) {

	newdhcpserver := r.FormValue("dhchserver")
	ticket := r.FormValue("ticket")

	if newdhcpserver == "" || ticket == "" {
		respondResonJSON(nil, w, r, "Dhcpserver or ticket field cannot be empty")
	} else {
		if utils.Exists("/edge/init") {
			data, err := ioutil.ReadFile("/edge/init")
			if err != nil {
				log.WithFields(log.Fields{"Operating": "migrate"}).Error(err)
			}
			deviceinfo := regeister.DeviceInfo{}
			yaml.Unmarshal(data, &deviceinfo)
			if deviceinfo.WorkID == "" {
				respondResonJSON(nil, w, r, "Missing workerid information")
			} else {

				deviceinfo.DhcpServer = newdhcpserver
				deviceinfo.Key = ""
				outdata, err := yaml.Marshal(deviceinfo)
				if err != nil {
					log.WithFields(log.Fields{"Operating": "migrate"}).Error(err)
				}
				f, err := os.OpenFile("/edge/init", os.O_WRONLY|os.O_TRUNC, 0777)
				defer f.Close()
				if err != nil {
					log.Error(err)
				}
				f.Write(outdata)
				log.WithFields(log.Fields{"Operating": "migrate"}).Info("Generate a device information file corresponding to the new cluster")
				regeister.BuildDeviceInfo(deviceinfo.Vpn, ticket, newdhcpserver)
				_, err = regeister.GetMyIP(regeister.WireGuardInterface)
				if err != nil {
					regeister.CloseWg()
				}
				_, err = regeister.GetMyIP(regeister.OpenVPNInterface)
				if err != nil {
					regeister.Closeopenvpn()
				}
			}

		} else {
			respondResonJSON(nil, w, r, "Equipment has not yet initialization")
		}

		//启动注册流程

		respondSuccessJSON("", w, r)
	}

}
