package system

import (
	"jxcore/internal/config"
	"jxcore/oplog"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"jxcore/web/controller/utils"
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
