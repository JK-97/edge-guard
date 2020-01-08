package network

import (
	"net/url"
	"os/exec"
)

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

// GetHost 从 url 中解析 Host
func GetHost(u string) string {
	uri, err := url.Parse(u)
	if err != nil {
		return u
	}
	return uri.Hostname()
}
