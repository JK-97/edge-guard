package updatemanage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/internal/network/dns"
	"jxcore/lowapi/system"
	"net/http"
	"os"
	"sync"
)

// 设置目标版本后，在后台自动更新。更新成功后上报版本，重启jxcore。
type UpdateManager interface {
	// 开始后台更新
	Start()

	// 获取当前更新状态
	GetStatus() UpdateStatus
	// 获取当前版本
	GetCurrentVersion() VersionInfo
	// 设置目标版本
	SetTargetVersion(version VersionInfo) error
	// 上传压缩包并更新
	UpdateWithZip(io.Reader) error
	// 上报版本号
	ReportVersion() error
}

var ErrUpdating = errors.New("updating")

var once *sync.Once = &sync.Once{}

func NewUpdateManager() UpdateManager {
	once.Do(func() {
		process = &updateManager{
			targetVersion:  parseTargetVersionFile(),
			currentVersion: parseCurrentVersionFile(),
			status:         FINISHED,

			targetVersionUpdated: make(chan bool),
		}
	})
	return process
}

func (up *updateManager) Start() {
	go up.updateLoop()
}

func (up *updateManager) GetStatus() UpdateStatus {
	return up.status
}

func (up *updateManager) GetCurrentVersion() VersionInfo {
	return up.currentVersion
}

func (up *updateManager) SetTargetVersion(version VersionInfo) error {
	if up.status == UPDATING {
		return ErrUpdating
	}

	up.targetVersion = version
	data, err := json.MarshalIndent(version, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(TARGETVERSION, data, 0644)
	if err != nil {
		return err
	}

	// send a non-block message to signal target version update
	select {
	case up.targetVersionUpdated <- true:
	default:
	}

	return nil
}

type Respdatastruct struct {
	Status   string            `json:"status"`
	WorkerId string            `json:"worker_id"`
	PkgInfo  map[string]string `json:"pkg_info"`
}

func (up *updateManager) ReportVersion() error {
	deviceinfo, _ := device.GetDevice()
	resprawinfo := Respdatastruct{
		Status:   string(up.GetStatus()),
		WorkerId: deviceinfo.WorkerID,
		PkgInfo:  up.GetCurrentVersion(),
	}
	respdata, err := json.Marshal(resprawinfo)
	if err != nil {
		return err
	}

	ip, port := dns.ParseIpInTxt(UPLOADDOMAIN)

	_, err = http.Post(fmt.Sprintf("http://%s:%s%s", ip, port, UPLOADPATH), "application/json", bytes.NewReader(respdata))
	if err != nil {
		return err
	}
	return nil
}

func (up *updateManager) UpdateWithZip(reader io.Reader) error {
	// TODO: add lock
	if up.status == UPDATING {
		return ErrUpdating
	}
	up.status = UPDATING
	defer func() { up.status = FINISHED }()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/tmp/update_package.tar", data, 0755)
	if err != nil {
		return err
	}
	_ = os.RemoveAll("/tmp/update_package")
	_ = os.MkdirAll("/tmp/update_package", 0666)
	err = system.RunCommand("tar -xvf /tmp/update_package.tar -C /tmp/update_package")
	if err != nil {
		return err
	}
	err = system.RunCommand("dpkg -i /tmp/update_package/debs/*.deb")
	if err != nil {
		return err
	}
	up.currentVersion = parseCurrentVersionFile()
	return nil
}
