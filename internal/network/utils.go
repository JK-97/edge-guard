package network

import "os/exec"

func CheckNetwork() bool {
	return Ping(testDomain)
}

func CheckMasterConnect() bool {
	return Ping(masterDomain)
}

func Ping(hostName string) bool {
	err := exec.Command("ping", hostName, "-c", "1", "-W", "5").Run()
	return err == nil
}
