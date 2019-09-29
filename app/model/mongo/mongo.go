package mongo

import (
	"context"
	"jxcore/app/model/clean"
	"jxcore/app/model/supervisor"
	"jxcore/config"
	"jxcore/log"
	"jxcore/utils"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var supervisorclient = supervisor.NewSupervisorRPC(config.InterSettings.Supervisor.Host)

///MongoCheck is
func MongoCheck() (err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	return err
}

//UnInstallMongo is
func UnInstallMongo() {
	supervisorclient.StopProcess(config.InterSettings.Mongodb.MongodbSupervisor)
	for _, pkgname := range config.InterSettings.Mongodb.MongoPkg {
		cmd := "dpkg -r " + pkgname
		uninstallcmd := exec.Command("/bin/sh", "-c", cmd)
		err := uninstallcmd.Start()
		if err != nil {
			log.Error(err)
		}
		uninstallcmd.Wait()
	}

	rccmd := exec.Command("/bin/sh", "-c", "dpkg -l | grep ^rc | cut -d' ' -f3 | sudo xargs dpkg --purge")
	rccmd.Start()
	if utils.Exists(config.InterSettings.Mongodb.MongoDataPath) {
		clean.DelFile([]string{config.InterSettings.Mongodb.MongoDataPath})
	}
}
func InstallMongo() {
	if utils.Exists(config.InterSettings.Restore.MongoPath) {
		installmongocmd := exec.Command("dpkg", "-i", "-R", config.InterSettings.Restore.MongoPath)
		out, err := installmongocmd.CombinedOutput()
		if err != nil {
			log.Error(err, string(out))
		} else {
			log.Info(string(out))
		}

		supervisorclient.StartProcess(config.InterSettings.Mongodb.MongodbSupervisor)
	}
}
