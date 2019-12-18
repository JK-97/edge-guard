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

	//固件升级
	r.HandleFunc("/edgenode/exec/update", controller.UpdateByDeb).Methods(http.MethodPost)

	v1Router := r.PathPrefix("/api/v1").Subrouter()
	// 登陆登出
	v1Router.HandleFunc("/login", controller.PostLogin).Methods(http.MethodPost)
	v1Router.HandleFunc("/logout", controller.PostLogout).Methods(http.MethodPost)

	secretRouter := v1Router.NewRoute().Subrouter()
	secretRouter.Use(requireLoginMiddleware)

	//操作
	secretRouter.HandleFunc("/system/upgrade", controller.UploadAndUpdate).Methods(http.MethodPost)
	secretRouter.HandleFunc("/system/reboot", controller.Reboot).Methods(http.MethodPost)

	//设置时间
	secretRouter.HandleFunc("/settings/time", controller.GetNtpConfig).Methods(http.MethodGet)
	secretRouter.HandleFunc("/settings/time", controller.SetNtpConfig).Methods(http.MethodPost)
	secretRouter.HandleFunc("/time", controller.GetTime).Methods(http.MethodGet)
	secretRouter.HandleFunc("/time", controller.SetTime).Methods(http.MethodPost)

	//更改信息
	secretRouter.HandleFunc("/node/info", controller.GetDeviceInfo).Methods(http.MethodGet)
	secretRouter.HandleFunc("/node/name", controller.SetDeviceName).Methods(http.MethodPost)

	// 驱动
	secretRouter.HandleFunc("/drivers", controller.GetEdgexDrivers).Methods(http.MethodGet)

	// 密码
	secretRouter.HandleFunc("/system/password", controller.SetPasswordHandler).Methods(http.MethodPost)

	return r
}
