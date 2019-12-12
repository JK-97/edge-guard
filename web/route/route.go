package route

import (
	"jxcore/web/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(logRequestMiddleware)
	r.Use(recoverMiddleware)

	//操作
	r.HandleFunc("/api/v1/system/upgrade", controller.UploadAndUpdate).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/system/reboot", controller.Reboot).Methods(http.MethodPost)

	//设置时间
	r.HandleFunc("/api/v1/settings/time", controller.GetNtpConfig).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/settings/time", controller.SetNtpConfig).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/time", controller.GetTime).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/time", controller.SetTime).Methods(http.MethodPost)
	//更改信息
	r.HandleFunc("/api/v1/node/info", controller.GetDeviceInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/node/name", controller.SetDeviceName).Methods(http.MethodPost)

	//固件升级
	r.HandleFunc("/edgenode/exec/update", controller.UpdateByDeb).Methods(http.MethodPost)

	return r
}
