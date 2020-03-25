package driver

import (
	"context"
	"jxcore/lowapi/docker"
	"jxcore/web/controller/utils"
	"jxcore/web/remote"
	"net/http"
	"time"
)

type GetEdgexDriversRespDriver struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProxyName   string `json:"proxyName"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Version     string `json:"version"`
	Alive       bool   `json:"alive"`
}

type GetEdgexDriversResp struct {
	Drivers []GetEdgexDriversRespDriver `json:"drivers"`
}

func GetEdgexDrivers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	deviceServices, err := remote.GetDeviceServices(ctx)
	if err != nil {
		panic(err)
	}

	resp := GetEdgexDriversResp{}
	for _, ds := range deviceServices {
		resp.Drivers = append(resp.Drivers, GetEdgexDriversRespDriver(*ds))
	}

	utils.RespondSuccessJSON(resp, w)
}

func isAuthedImage(fileName string) (name, version string, err error) {
	return
}

func PostInstallDriver(w http.ResponseWriter, r *http.Request) {
	file, fileInfo, err := r.FormFile("image")
	if err != nil {
		panic(err)
	}
	_, _, err = isAuthedImage(fileInfo.Filename)
	if err != nil {
		panic(err)
	}
	err = docker.LoadImage(file)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(nil, w)
}
