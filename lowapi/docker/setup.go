package docker

import (
	"encoding/json"
	"github.com/JK-97/edge-guard/internal/network"
	"io/ioutil"
	"os/exec"

	"github.com/JK-97/edge-guard/lowapi/logger"
)

func EnsureDockerDNSConfig() error {
	logger.Info("Checking docker DNS config")
	// 解析daemon.json 文件
	b, err := ioutil.ReadFile(daemonConfigPath)
	if err != nil {
		return err
	}
	var conf map[string]interface{}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}

	// 如果dns没有配置，添加配置，并重启docker
	if !addDnsConf(conf) {
		return nil
	}
	out, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(daemonConfigPath, out, 0666)
	if err != nil {
		return err
	}

	return exec.Command("systemctl", "restart", "docker").Run()
}

func addDnsConf(conf map[string]interface{}) (needRestart bool) {
	dnsConf, ok := conf["dns"]
	if !ok {
		conf["dns"] = []string{network.DockerHostIP}
		return true
	}
	l, ok := dnsConf.([]interface{})
	if !ok {
		conf["dns"] = []string{network.DockerHostIP}
		return true
	}

	for _, d := range l {
		if s, ok := d.(string); ok && s == network.DockerHostIP {
			return false
		}
	}
	conf["dns"] = append(l, network.DockerHostIP)
	return true
}
