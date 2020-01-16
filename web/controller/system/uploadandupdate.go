package system

import (
	"jxcore/lowapi/system"
	"jxcore/management/updatemanage"
	"jxcore/web/controller/utils"
	"net/http"
	"time"
)

func UploadAndUpdate(w http.ResponseWriter, r *http.Request) {
	manager := updatemanage.NewUpdateManager()
	err := manager.UpdateWithZip(r.Body)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
	system.RestartJxcoreAfter(5 * time.Second)
}
