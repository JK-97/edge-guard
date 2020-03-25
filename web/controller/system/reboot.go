package system

import (
	"jxcore/lowapi/system"
	"jxcore/oplog"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"jxcore/web/controller/utils"
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
