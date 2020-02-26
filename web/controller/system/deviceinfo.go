package system

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/lowapi/store/filestore"
	"jxcore/management/updatemanage"
	"jxcore/web/controller/utils"
	"net/http"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

const (
	deviceNameKey = "name"
)

type deviceInfoRequest struct {
	Name             string `json:"name"`
	SkipUploadMaster bool   `json:"skip_upload_master"`
}

type deviceinfoResponse struct {
	Name            string `json:"name"`
	WorkerID        string `json:"workerid"`
	Model           string `json:"model"`
	FirmwareVersion string `json:"firmware_version"`
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	currentDevice, err := device.GetDevice()
	if err != nil {
		panic(err)
	}
	deviceType, err := device.GetDeviceType()
	if err != nil {
		panic(err)
	}

	firmwareVersion := updatemanage.NewUpdateManager().GetCurrentVersion()

	deviceNameData, err := filestore.KV.GetDefault(deviceNameKey, []byte(currentDevice.WorkerID))
	if err != nil {
		logger.Error("failed to get device name", err)
	}

	reponse := &deviceinfoResponse{
		Name:     string(deviceNameData),
		WorkerID: currentDevice.WorkerID,
		// TODO 需要产品定义
		Model:           string(deviceType),
		FirmwareVersion: firmwareVersion["jx-toolset"],
	}

	utils.RespondJSON(reponse, w, 200)
}

// 设置设备名称
func SetDeviceName(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	request := deviceInfoRequest{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		panic(err)
	}

	err = filestore.KV.Set(deviceNameKey, []byte(request.Name))
	if err != nil {
		panic(err)
	}
	if !request.SkipUploadMaster {
		reportMaster(request.Name)
	}
	utils.RespondSuccessJSON(nil, w)
}

// 上报云端
func reportMaster(name string) {
	// TODO
}
