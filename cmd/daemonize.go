package cmd

import (
	daemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

func Deamonize(proc func()) {
	context := daemon.Context{
		PidFileName: "/var/run/edge-guard.pid",
		PidFilePerm: 0644,
		LogFileName: "/edge/logs/edge-guard.log",
		LogFilePerm: 0644,
		// LogFileName: "/dev/stdout"
	}

	child, err := context.Reborn()
	if err != nil {
		context := daemon.Context{
			PidFileName: "/var/run/edge-guard.pid",
			PidFilePerm: 0644,
			LogFileName: "/edge/logs/edge-guard.log",
			LogFilePerm: 0644,
		}
		child, err = context.Reborn()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatal("Unable to run")
		}
	}
	if child != nil {
		return
	}
	defer context.Release()
	proc()
}
