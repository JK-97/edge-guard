package controller

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/core"
	"jxcore/core/device"
	log "jxcore/go-utils/logger"
	"jxcore/lowapi/utils"
	"jxcore/management/updatemanage"
	"net/http"
)

func UpdateByDeb(w http.ResponseWriter, r *http.Request) {
	core := core.GetJxCore()
	updateprocess := updatemanage.GetUpdateProcess()
	if updateprocess.GetStatus() != updatemanage.FINISHED {
		w.WriteHeader(400)
		respondResonJSON(nil, w, r, "machine is busy to updating,please update later")
	} else {
		reqrawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
		}
		reqinfo := updatemanage.Reqdatastruct{}
		err = json.Unmarshal(reqrawdata, &reqinfo)
		if err != nil {
			log.Error(err)
		}
		indentdata, err := json.MarshalIndent(reqinfo.Data, "", "  ")
		if err != nil {
			log.Error(err)
			w.WriteHeader(400)
			respondResonJSON(nil, w, r, "json format err")
		} else {
			updateprocess.SetNewTarget(indentdata)
			go func() {
				core.UpdateCore()
			}()
			deviceinfo, err := device.GetDevice()
			utils.CheckErr(err)
			respdata := updatemanage.Respdatastruct{
				Status:   updatemanage.UPDATING.String(),
				WorkerId: deviceinfo.WorkerID,
				PkgInfo:  updatemanage.ParseVersionFile(),
			}
			respondJSON(respdata, w, r)
			//wg.Add(2)

		}

	}

}
