package system

import (
	"fmt"
	"jxcore/lowapi/system"
	"jxcore/management/updatemanage"
	"jxcore/oplog"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"jxcore/web/controller/utils"
	"net/http"
	"time"
)

func UploadAndUpdate(w http.ResponseWriter, r *http.Request) {
	manager := updatemanage.NewUpdateManager()
	err := manager.UpdateWithZip(r.Body)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
	oplog.Insert(logs.NewOplog(types.UPDATE, fmt.Sprintf("updated by upload")))
	system.RestartJxcoreAfter(5 * time.Second)
}
