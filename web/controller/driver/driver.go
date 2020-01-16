package driver

import (
	"context"
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
