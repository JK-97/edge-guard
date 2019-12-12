package controller

import (
	"net/http"
)

const (
	urlGetDeviceServices = "edgegw.iotedge:48081/api/v1/deviceservice"
)

func GetEdgexDrivers(w http.ResponseWriter, r *http.Request) {
	http.Get("")
}
