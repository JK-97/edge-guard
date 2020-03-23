package route

import (
	"jxcore/web/controller/driver"
	"jxcore/web/controller/system"
	"net/http"

	"github.com/gorilla/mux"
)

const staticFilePath = "/edge/jxcore/frontend"

func Routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(logRequestMiddleware)
	r.Use(recoverMiddleware)

	//固件升级
	r.HandleFunc("/edgenode/exec/update", system.UpdateByDeb).Methods(http.MethodPost)

	v1Router := r.PathPrefix("/api/v1").Subrouter()
	v1Router.HandleFunc("/log", system.GetOplog).Methods(http.MethodGet)
	v1Router.HandleFunc("/log/download", system.DownloadOplog).Methods(http.MethodGet)
	// 登陆登出
	v1Router.HandleFunc("/login", system.PostLogin).Methods(http.MethodPost)
	v1Router.HandleFunc("/logout", system.PostLogout).Methods(http.MethodPost)

	secretRouter := v1Router.NewRoute().Subrouter()
	secretRouter.Use(requireLoginMiddleware)

	//操作
	secretRouter.HandleFunc("/system/upgrade", system.UploadAndUpdate).Methods(http.MethodPost)
	secretRouter.HandleFunc("/system/reboot", system.Reboot).Methods(http.MethodPost)

	//设置时间
	secretRouter.HandleFunc("/settings/time", system.GetNtpConfig).Methods(http.MethodGet)
	secretRouter.HandleFunc("/settings/time", system.SetNtpConfig).Methods(http.MethodPost)
	secretRouter.HandleFunc("/time", system.GetTime).Methods(http.MethodGet)
	secretRouter.HandleFunc("/time", system.SetTime).Methods(http.MethodPost)

	//更改信息
	secretRouter.HandleFunc("/node/info", system.GetDeviceInfo).Methods(http.MethodGet)
	secretRouter.HandleFunc("/node/name", system.SetDeviceName).Methods(http.MethodPost)

	// 驱动
	secretRouter.HandleFunc("/drivers", driver.GetEdgexDrivers).Methods(http.MethodGet)
	secretRouter.HandleFunc("/drivers", driver.PostInstallDriver).Methods(http.MethodPost)

	// 代理请求到device service
	secretRouter.HandleFunc("/driver/{dsname}{path:.*}", driver.Proxy)

	// 网络
	secretRouter.HandleFunc("/network/interfaces", system.GetNetworkInterfaces).Methods(http.MethodGet)
	secretRouter.HandleFunc("/network/interface/{iface}", system.GetNetworkInterfaceByName).Methods(http.MethodGet)
	secretRouter.HandleFunc("/network/interfaces/fourg", system.GetFourGInterface).Methods(http.MethodGet)
	secretRouter.HandleFunc("/network/interfaces/fourg", system.EnableFourGInterface).Methods(http.MethodPost)

	//日志系统
	// secretRouter.HandleFunc("/log", system.GetOplog).Methods(http.MethodGet)
	// secretRouter.HandleFunc("/log/download", system.DownloadOplog).Methods(http.MethodGet)

	// 密码
	secretRouter.HandleFunc("/system/password", system.SetPasswordHandler).Methods(http.MethodPost)

	// 前端
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticFilePath)))

	return r
}
