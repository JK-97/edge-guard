package system

import (
	"jxcore/internal/config"
	"jxcore/web/controller/utils"
	"net/http"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	err := config.ResetSystemConfig()
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
}
