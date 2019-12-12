package controller

import (
	"io/ioutil"
	"jxcore/lowapi/system"
	"jxcore/management/updatemanage"
	"net/http"
)

func UploadAndUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("/tmp/update_package.tar", data, 0755)
	if err != nil {
		panic(err)
	}
	err = system.RunCommand("tar -xvf /tmp/update_package.tar -o /tmp/update_package")
	if err != nil {
		panic(err)
	}
	err = system.RunCommand("dpkg -i /tmp/update_package/*")
	if err != nil {
		panic(err)
	}
	updateprocess := updatemanage.GetUpdateProcess()
	updateprocess.ReportVersion()
	RespondSuccessJSON(nil, w)
	system.RestartJxcoreAfter(0)
}
