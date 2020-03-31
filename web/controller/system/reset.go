package system

import (
	"github.com/JK-97/edge-guard/internal/config"
	"github.com/JK-97/edge-guard/oplog"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"github.com/JK-97/edge-guard/web/controller/utils"
	"net/http"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	oplog.Insert(logs.NewOplog(types.DEVICE, "reset system"))
	err := config.ResetSystemConfig()
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
}
