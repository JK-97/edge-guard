package config

import (
	log "jxcore/lowapi/logger"
	"os"
	"text/template"
)

const (
	tmplDir = "/edge/jxcore/template"

	telegrafTmplPath = tmplDir + "/telegraf.cfg.tpl"
	telegrafConfPath = "/edge/monitor/telegraf/bin/telegraf.cfg"

	statsiteTmplPath = tmplDir + "/influxdb.ini.tpl"
	statsiteConfPath = "/edge/monitor/device-statsite/conf/influxdb.ini"

	cadvisorTmplPath = tmplDir + "/cadvisor.yaml.tpl"
	cadvisorConfPath = "/jxbootstrap/worker/docker-compose.d/cadvisor/docker-compose.yaml"

	timesyncdTmplPath = tmplDir + "/timesyncd.conf.tpl"
	timesyncdConfPath = "/etc/systemd/timesyncd.conf"
)

func Telegrafcfg(masterip string, workerid string) {
	type MASTER struct {
		WORKER_ID string
		MASTER_IP string
	}
	sweaters := MASTER{MASTER_IP: masterip, WORKER_ID: workerid}
	err := putConfig(telegrafTmplPath, telegrafConfPath, sweaters)
	if err != nil {
		log.Error(err)
	}
}

func Statsitecfg(masterip string, vpnIP string) {
	type MASTER struct {
		MASTER_IP string
		VpnIP     string
	}
	sweaters := MASTER{MASTER_IP: masterip, VpnIP: vpnIP}
	err := putConfig(statsiteTmplPath, statsiteConfPath, sweaters)
	if err != nil {
		log.Error(err)
	}
}

func Cadvisorcfg(masterip string, workerid string) {
	type MASTER struct {
		MASTER_IP string
		WORKER_ID string
	}
	sweaters := MASTER{MASTER_IP: masterip, WORKER_ID: workerid}
	err := putConfig(cadvisorTmplPath, cadvisorConfPath, sweaters)
	if err != nil {
		log.Error(err)
	}
}

func TimdsyncdCfg(serverAddr string) {
	type MASTER struct {
		SERVER_ADDR string
	}
	sweaters := MASTER{SERVER_ADDR: serverAddr}
	err := putConfig(timesyncdTmplPath, timesyncdConfPath, sweaters)
	if err != nil {
		log.Error(err)
	}

}

func putConfig(tmplPath, targetPath string, sweaters interface{}) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return tmpl.Execute(f, sweaters)
}
