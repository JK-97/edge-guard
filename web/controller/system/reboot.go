package system

import (
	"jxcore/lowapi/system"
	"jxcore/web/controller/utils"
	"net/http"
)

func Reboot(w http.ResponseWriter, r *http.Request) {
	err := system.RebootAfter(0)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
}
