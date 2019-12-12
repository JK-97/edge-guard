package controller

import (
	"jxcore/lowapi/system"
	"net/http"
)

func Reboot(w http.ResponseWriter, r *http.Request) {
	err := system.RebootAfter(0)
	if err != nil {
		panic(err)
	}
	RespondSuccessJSON(nil, w)
}
