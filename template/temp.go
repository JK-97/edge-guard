package template

import (
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"os"
	"text/template"
)

func Telegrafcfg(masterip string, workerid string) {
	tmpl, err := template.ParseFiles("/edge/jxcore/template/telegraf.cfg.tpl")
	if err != nil {
		log.Error(err)
	}
	type MASTER struct {
		WORKER_ID string
		MASTER_IP string
	}
	sweaters := MASTER{MASTER_IP: masterip, WORKER_ID: workerid}
	//name :=masterip
	f, err := os.OpenFile("/edge/monitor/telegraf/bin/telegraf.cfg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err)
	}
	err = tmpl.Execute(f, sweaters)
	if err != nil {
		log.Error(err)
	}
}

func Statsitecfg(masterip string, vpnIP string) {
	tmpl, err := template.ParseFiles("/edge/jxcore/template/influxdb.ini.tpl")
	if err != nil {
		log.Error(err)
	}
	type MASTER struct {
		MASTER_IP string
		VpnIP     string
	}
	sweaters := MASTER{MASTER_IP: masterip, VpnIP: vpnIP}
	//name :=masterip
	f, err := os.OpenFile("/edge/monitor/device-statsite/conf/influxdb.ini", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
	}
	err = tmpl.Execute(f, sweaters)
	if err != nil {
		log.Error(err)
	}
}

func Cadvisorcfg(masterip string, workerid string) {
	tmpl, err := template.ParseFiles("/edge/jxcore/template/cadvisor.yaml.tpl")
	if err != nil {
		log.Error(err)
	}
	type MASTER struct {
		MASTER_IP string
		WORKER_ID string
	}
	sweaters := MASTER{MASTER_IP: masterip, WORKER_ID: workerid}
	//name :=masterip
	f, err := os.OpenFile("/jxbootstrap/worker/docker-compose.d/cadvisor/docker-compose.yaml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
	}
	err = tmpl.Execute(f, sweaters)
	if err != nil {
		log.Error(err)
	}
}
