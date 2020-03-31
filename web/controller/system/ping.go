package system

import (
	"github.com/JK-97/edge-guard/web/controller/utils"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON("pong", w, 200)
}
