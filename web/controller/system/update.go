package system

import (
	"jxcore/core/device"
	"jxcore/management/updatemanage"
	"jxcore/web/controller/utils"
	"net/http"
)

type reqdatastruct struct {
	Data map[string]string `json:"data"`
}

func UpdateByDeb(w http.ResponseWriter, r *http.Request) {
	manager := updatemanage.NewUpdateManager()

	reqinfo := reqdatastruct{}
	utils.MustUnmarshalJson(r.Body, &reqinfo)

	err := manager.SetTargetVersion(reqinfo.Data)
	if err == updatemanage.ErrUpdating {
		utils.RespondReasonJSON(nil, w, "machine is busy to updating, please update later", 400)
		return
	}
	if err != nil {
		panic(err)
	}

	deviceinfo, err := device.GetDevice()
	if err != nil {
		panic(err)
	}
	respdata := updatemanage.Respdatastruct{
		Status:   string(updatemanage.UPDATING),
		WorkerId: deviceinfo.WorkerID,
		PkgInfo:  manager.GetCurrentVersion(),
	}
	utils.RespondJSON(respdata, w, 200)
}
