package network

import (
	"fmt"
	log "jxcore/go-utils/logger"
	"os/exec"
	"time"
)

const (
	checkBestIFaceInterval = time.Second * 5
	testServer             = "114.114.114.114"
)

var (
	currentIFace  string
	ifacePriority = []string{"eth0", "eth1", "usb0", "usb1"}
)

func findBestIFace() string {
	for _, iface := range ifacePriority {
		connected, err := testConnect(iface, testServer)
		if err != nil {
			log.Error("Test connect err:", err)
		}
		if connected {
			return iface
		}
	}
	return ifacePriority[len(ifacePriority)-1]
}

func switchIFace(iFace string) (err error) {
	err = exec.Command("ifdown", currentIFace).Run()
	if err != nil {
		return
	}
	err = exec.Command("ifup", iFace, "--force").Run()
	if err != nil {
		return
	}
	currentIFace = iFace
	return
}

func InitIFace() error {
	return switchIFace(findBestIFace())
}

func MaintainBestIFace() error {
	timer := time.NewTicker(checkBestIFaceInterval)
	defer timer.Stop()

	for range timer.C {
		bestIFace := findBestIFace()
		if currentIFace != bestIFace {
			err := switchIFace(bestIFace)
			log.Error("Failed to switch network interface", err)
		}
	}
	return nil
}

func setIPRoute(netInterface string) (err error) {
	// err = exec.Command("ip", "route", "add", "114.114.114.114/32", "via", netRoute, "dev", netInterface).Run()
	err = exec.Command("/bin/bash", "-c", fmt.Sprintf("ip route replace %s/32 dev %s", testServer, netInterface)).Run()
	log.Info("setIPRoute:", netInterface+"114.114.114.114/32 router")
	// waitForSet(netInterface)
	return
}

func removeIPRoute(netInterface string) (err error) {
	err = exec.Command("/bin/bash", "-c", fmt.Sprintf("ip route del %s/32 dev %s", testServer, netInterface)).Run()
	log.Info("removeIPRoute :", netInterface+"114.114.114.114/32 router")
	// waitForRemove(netInterface)
	return
}

func testConnect(netInterface, dhcpHost string) (connected bool, err error) {
	err = setIPRoute(netInterface)
	if err != nil {
		return
	}
	pingErr := exec.Command("ping", "-c", "1", dhcpHost).Run()

	if err != nil {
		if exitError, ok := pingErr.(*exec.ExitError); ok {
			log.Infof("Ping from interface %v to %v exited with code %v", netInterface, dhcpHost, exitError.ExitCode())
			connected = false
		} else {
			err = pingErr
		}
	} else {
		connected = true
		log.Info("checkRoute OK : ", netInterface)
	}

	// err = removeIPRoute(netRoute, netInterface)

	return
}
