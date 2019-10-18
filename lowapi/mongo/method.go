package mongo

import (
    "context"
    log "jxcore/go-utils/logger"
    "jxcore/lowapi/supervisor"
    "jxcore/lowapi/utils"
    "os/exec"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var supervisorclient = supervisor.NewSupervisorRPC(supervisor.SupervisorHost)

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
   supervisorclient.StopProcess("")
   for _, pkgname := range MongooPKg{
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
   if utils.Exists(MongoDataPath) {
       utils.DelFile([]string{MongoDataPath})
   }
}
func InstallMongo() {
   if utils.Exists(MongoRestore) {
       installmongocmd := exec.Command("dpkg", "-i", "-R", MongoRestore)
       out, err := installmongocmd.CombinedOutput()
       if err != nil {
           log.Error(err, string(out))
       } else {
           log.Info(string(out))
       }

       supervisorclient.StartProcess( "")
   }
}
