/*
	jxcore 高阶api
*/
package controller

import (
	"jxcore/app/model/clean"
	"jxcore/app/model/docker"
	"jxcore/app/model/mongo"
	"jxcore/app/model/pythonpackage"
	"jxcore/config"
	"net/http"
	"jxcore/log"
)

var dockerobj = docker.NewClient()
var c = pythonpackage.NewPkgClient()

//从恢复目录恢复 docker python
//清除mongo中的配置与数据，为迁移注册作准备
func Restore(w http.ResponseWriter, r *http.Request) {
	//恢复docker
	dockerobj.ContainerAllRemove()
	dockerobj.ImageAllRemove()
	go dockerobj.DockerRestore()
	//恢复 python
	err := c.DeletePyPkg()
	if err != nil {
		log.Error(err)
	}
	go c.RestorePyPkg()
	//清除恢复 mongo
	mongo.UnInstallMongo()
	mongo.InstallMongo()
	// 删除一些必要文件
	clean.DelFile(config.InterSettings.DelPath)
	clean.ResetFile(config.InterSettings.CleanPath)
	respondSuccessJSON("", w, r)

}
