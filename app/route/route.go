package route

import (
	"jxcore/app/controller"
	"net/http"
)

// Routes adds routes to http
func Routes() http.Handler {
	mux := http.NewServeMux()
	handler := logRequest(mux)

	mux.HandleFunc("/api/v1/ping", controller.PingGet)

	mux.HandleFunc("/", controller.IndexGET)

	//docker
	//mux.HandleFunc("/docker/info/images", controller.DockerImagesGET)
	//mux.HandleFunc("/docker/info/container", controller.DockerContainerGET)
	//mux.HandleFunc("/docker/delete", controller.DockerRemoveDelete)
	//mux.HandleFunc("/docker/exec/restore", controller.DockerRestorePost)
	//
	////monogo
	//mux.HandleFunc("/mongo/delete", controller.MongoRemoveDelete)
	////supervisor
	//mux.HandleFunc("/supervisor/info/all", controller.SupervisorAllProcessGET)
	//mux.HandleFunc("/supervisor/exec/stop", controller.SupervisorRestoreProcessPost)
	//
	////pythonpkg
	//mux.HandleFunc("/pythonpkg/info/all", controller.AllPythonPkgInfoGet)
	//mux.HandleFunc("/pythonpkg/info/internal", controller.InternalPythonPkgInfoGet)
	//mux.HandleFunc("/pythonpkg/delete", controller.PythonPkgUinstallDelete)
	//mux.HandleFunc("/pythonpkg/exec/restore", controller.PythonPkgURestorePost)
	//test √
	mux.HandleFunc("/edgenode/exec/restore", controller.Restore)
	//test √
	mux.HandleFunc("/edgenode/exec/clean", controller.CleanDelete)
	//test √
	mux.HandleFunc("/edgenode/exec/update", controller.Update)
	mux.HandleFunc("/edgenode/exec/migrate", controller.Migrate)
	//test √
	mux.HandleFunc("/edgenode/exec/reload", controller.Reload)
	//test √
	mux.HandleFunc("/edgenode/changelog", controller.ChangeLog)
	//test √
	mux.HandleFunc("/edgenode/version", controller.Version)

	mux.HandleFunc("/edgenode/componentstate", controller.ComponentState)

	return handler
}
