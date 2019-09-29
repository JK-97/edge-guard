/*
	低阶api
 */
package v1

import (
	"github.com/pkg/errors"
	"jxcore/app/model/docker"
	"net/http"
)

var dockerobj = docker.NewClient()

//DockerImagesGET Http GET
//获取所有的docker 镜像
func DockerImagesGET(w http.ResponseWriter, r *http.Request) {
	resp, err := dockerobj.ImagesList()
	if err != nil {
		log.Error(err)
	}
	respondSuccessJSON(resp, w, r)
}

//DockerContainerGET Http GET
//获取所有的docker容器
func DockerContainerGET(w http.ResponseWriter, r *http.Request) {
	resp, err := dockerobj.ContainerList()
	if err != nil {
		log.Error(err)
	}
	respondSuccessJSON(resp, w, r)
}

//DockerRemoveDelete Http DELETE
//删除所有数据
func DockerRemoveDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {

		go dockerobj.ContainerAllRemove()
		go dockerobj.ImageAllRemove()
		respondSuccessJSON(nil, w, r)
	} else {
		Error(w, errors.New("meathod not support"), 404)
	}
}

//DockerRestorePost Http Post
//从还原文件中中恢复镜像
func DockerRestorePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		dockerobj.DockerRestore()
		respondSuccessJSON(nil, w, r)
	} else {
		Error(w, errors.New("meathod not support"), 404)
	}

}
