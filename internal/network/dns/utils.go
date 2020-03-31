package dns

import (
	"net"
	"os/exec"
	"strings"

	"github.com/JK-97/edge-guard/lowapi/logger"
	log "github.com/JK-97/edge-guard/lowapi/logger"
)

func RestartDnsmasq() {
	err := exec.Command("/bin/bash", "-c", "systemctl restart dnsmasq").Run()
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Restarted dnsmasq")
}

func ReloadDnsmasq() {
	err := exec.Command("/bin/bash", "-c", "systemctl reload dnsmasq").Run()
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Reloaded dnsmasq")
}

func ParseIpInTxt(url string) (string, string) {
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		log.Error("Possible DNS configuration error: ", err)
	}
	if len(txtRecords) == 0 {
		return "", ""
	}
	res := strings.Split(txtRecords[0], ":")
	return res[0], res[1]
}
