package system

import (
	"jxcore/web/controller/utils"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON("pong", w, 200)
}
