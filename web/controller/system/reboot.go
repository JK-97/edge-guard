package system

import (
	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/oplog"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"github.com/JK-97/edge-guard/web/controller/utils"
	"net/http"
)

func Reboot(w http.ResponseWriter, r *http.Request) {
	oplog.Insert(logs.NewOplog(types.DEVICE, "rebbot"))
	err := system.RebootAfter(0)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
}
