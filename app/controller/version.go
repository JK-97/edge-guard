package controller

import (
	"jxcore/app/model/version"
	"net/http"
)
//获取当前的组件版本信息
//	response:
//		{
//		"data": {
//		"camera": "2.0",
//		"edge": "1.2",
//		"rs485": "2.0",
//		"telegraf": "1.3"
//		},
//		"desc": "success"
//		}
func Version(w http.ResponseWriter, r *http.Request) {
	versioninfo := version.PraseVersionFile()
	respondSuccessJSON(versioninfo, w, r)
}


//获取组件的更新日志
//每个组件下有配置文件，修改是会主动触发changelog的更新
//version文件以内以yaml文件格式编写
//	eg:
//		name： edge
//		version: 1.2
//
//	response:
//		{
//		"2019-09-09 09:51:07": {
//		"camera": {
//		"2.1": "2019-09-09 09:52:41"
//		},
//		"edge": {
//		"1.3": "2019-09-09 09:51:07"
//		},
//		"rs485": {
//		"2.1": "2019-09-09 09:52:07"
//		},
//		"telegraf": {
//		"1.4": "2019-09-09 09:51:23",
//		"1.5": "2019-09-09 09:51:39"
//		}
//		},
//		"2019-09-09 09:53:41": {
//		"camera": {
//		"2.1": "2019-09-09 09:53:41"
//		},
//		"edge": {
//		"1.4": "2019-09-09 09:53:41"
//		},
//		"rs485": {
//		"2.1": "2019-09-09 09:53:41",
//		"2.2": "2019-09-09 09:56:20"
//		},
//		"telegraf": {
//		"1.5": "2019-09-09 09:53:41",
//		"1.6": "2019-09-09 09:55:28"
//		}
//		}
//		}
func ChangeLog(w http.ResponseWriter, r *http.Request) {
	changeloginfo :=version.ChangLog(version.PraseVersionFile())

	respondSuccessJSON(changeloginfo, w, r)
}
