package system

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/management/updatemanage"
	"github.com/JK-97/edge-guard/oplog"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"github.com/JK-97/edge-guard/web/controller/utils"
)

func UploadAndUpdate(w http.ResponseWriter, r *http.Request) {
	manager := updatemanage.NewUpdateManager()
	err := manager.UpdateWithZip(r.Body)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
	oplog.Insert(logs.NewOplog(types.UPDATE, fmt.Sprintf("updated by upload")))
	system.RestartEdgeguardAfter(5 * time.Second)
}
