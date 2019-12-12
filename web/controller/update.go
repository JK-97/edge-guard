package controller

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/core"
	"jxcore/core/device"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/utils"
	"jxcore/management/updatemanage"
	"net/http"
)

func UpdateByDeb(w http.ResponseWriter, r *http.Request) {
	updateprocess := updatemanage.GetUpdateProcess()
	if updateprocess.GetStatus() != updatemanage.FINISHED {
		RespondReasonJSON(nil, w, "machine is busy to updating,please update later", 400)
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
			RespondReasonJSON(nil, w, "json format err", 400)
		} else {
			updateprocess.SetNewTarget(indentdata)
			go func() {
				core.CheckCoreUpdate()
			}()
			deviceinfo, err := device.GetDevice()
			utils.CheckErr(err)
			respdata := updatemanage.Respdatastruct{
				Status:   updatemanage.UPDATING.String(),
				WorkerId: deviceinfo.WorkerID,
				PkgInfo:  updatemanage.ParseVersionFile(),
			}
			RespondJSON(respdata, w, 200)
		}

	}

}
