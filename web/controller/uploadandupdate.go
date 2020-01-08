package controller

import (
	"jxcore/lowapi/system"
	"jxcore/management/updatemanage"
	"net/http"
	"time"
)

func UploadAndUpdate(w http.ResponseWriter, r *http.Request) {
	manager := updatemanage.NewUpdateManager()
	err := manager.UpdateWithZip(r.Body)
	if err != nil {
		panic(err)
	}
	RespondSuccessJSON(nil, w)
	system.RestartJxcoreAfter(5 * time.Second)
}
