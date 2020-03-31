package system

import (
	"fmt"
	"github.com/JK-97/edge-guard/core/device"
	"github.com/JK-97/edge-guard/management/updatemanage"
	"github.com/JK-97/edge-guard/oplog"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"github.com/JK-97/edge-guard/web/controller/utils"
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
	oplog.Insert(logs.NewOplog(types.UPDATE, fmt.Sprintf("updated online %s", reqinfo.Data["jx-toolset"])))
	utils.RespondJSON(respdata, w, 200)
}
