package controller

import (
	"jxcore/component"
	"jxcore/log"
	"jxcore/utils"
	"net/http"
	"strings"
	"time"
)

//升级component
//	request：
//		file：更新的zip压缩宝文件
//		component：更新组件的名称
func Update(w http.ResponseWriter, r *http.Request) {

	//更新组件名字
	componenttoupdate := strings.ToLower(r.FormValue("component"))
	//更新的zip文件
	componentfile, _, err := r.FormFile("file")
	if err != nil {
		log.Debug(err)
		return
	}
	buff := make([]byte, 0)
	_, err = componentfile.Read(buff)
	if err != nil {
		log.Error(err)
	}
	defer componentfile.Close()

	component.StopComponent(componenttoupdate)
	//zip写入tmp
	//tempfilename := strconv.FormatInt(time.Now().UnixNano(), 10)
	//err=utils.SaveFile(tempfilename, componentfile)
	if err != nil {
		log.Error(err)
	}

	//解压到相应目录
	updatepath := component.ComponentPath[componenttoupdate]

	log.Error(updatepath)
	for {
		err = utils.Unzip(buff, updatepath)
		if err != nil {
			log.Error(err)
			time.Sleep(300 * time.Millisecond)
		} else {
			break
		}
	}

	//os.Remove(tempfilename)
	component.StartComponent(componenttoupdate)
	respondSuccessJSON("", w, r)

}
