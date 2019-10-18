package power

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ShutdownScript 关机脚本的位置
var ShutdownScript = "/app/bootstrap/scripts/shutdown.sh"

// getBattery 获取电池电量信息
func getBattery() int {
	if batteryEventPath == "" {
		batteryEventPath = "/sys/class/power_supply/bq3060-bat"
	}

	file, err := os.Open(batteryEventPath)
	if err != nil {
		return 0
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "POWER_SUPPLY_CAPACITY") {
			capacity, err := strconv.Atoi(strings.Split(s, "=")[1])
			if err != nil {
				return 0
			}
			return capacity
		}
	}

	return 0
}

// SystemPowerOff 关机
func SystemPowerOff(wakeTime int) {

	wakeTime *= 2

	log.Printf("shutting down. wake up after %d seconds", wakeTime)

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, ShutdownScript, strconv.Itoa(wakeTime))
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	} else {
		defer cmd.Wait()
	}
}

// StartUpMode 设备运行模式
// [休眠与唤醒](http://zentao.jiangxingai.com/zentao/doc-view-106.html)
type StartUpMode int

const (
	// ModeUnkonwn 未知模式
	ModeUnkonwn StartUpMode = iota

	// ModeManually 手动拍照：手动唤醒,手动拍一张照片, 然后进入休眠, 优先级1
	ModeManually

	// ModeManuallyWithTick 手动拍照10min：手动唤醒, 拍照异常, 10min后关机，优先级2
	ModeManuallyWithTick

	// ModeAuto 自动拍照：自动开机, 拍一张照片进入休眠，优先级3
	ModeAuto

	// ModeAutoWithTick 自动拍照10min：自动开机, 拍照片异常, 进入休眠 ， 优先级4
	ModeAutoWithTick
)

func getStartUpMode() StartUpMode {
	file, err := os.Open(startUpModePath)
	if err != nil {
		return 0
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b := scanner.Bytes()[0]
		return StartUpMode(int(b) - int('0'))
	}
	// TODO: 默认模式
	return ModeManually
}

// // PowerHandler 处理关机相关事项
// type PowerHandler struct {
// }
