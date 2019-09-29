package controller

import (
	"jxcore/log"
	"jxcore/monitor"
	"net/http"
	"os"
)

//杀死gateway的进程以及所有组件进程以加载新的配置文件
func Reload(w http.ResponseWriter, r *http.Request) {
	log.Error(monitor.GW.GwCMD.Process.Pid)
	//使用pid找到gateway 进行并杀死
	process,err:=os.FindProcess(monitor.GW.GwCMD.Process.Pid)
	if err !=nil{
		log.Error(err)
	}
	process.Kill()
	respondSuccessJSON("",w,r)

}
