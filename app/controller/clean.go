package controller

import (
	"jxcore/app/model/clean"
	"jxcore/config"
	"net/http"
)


//删除delpath下的所有文件，保存文件夹的格式
//重置cleanpath下的所有文件，保存文件的结构
func CleanDelete(w http.ResponseWriter, r *http.Request) {
	clean.DelFile(config.InterSettings.DelPath)
	clean.ResetFile(config.InterSettings.CleanPath)
	respondSuccessJSON(nil, w, r)
}
