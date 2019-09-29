package monitor

import (
	"io/ioutil"
	"jxcore/app/model/version"
	"jxcore/app/schema"
	"jxcore/log"
	"jxcore/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	"golang.org/x/sys/unix"
)

// DirSize is
func DirSize(path string) (int, error) {
	var size int
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += int(info.Size())
		}
		return err
	})
	return size, err
}

func SizeControl(cleanstrategy []schema.CleanStrategy, cleantimestep int) {
	for _, perfile := range cleanstrategy {
		cursize, _ := DirSize(perfile.Path)
		if cursize > perfile.Size {
			log.Warn("detect path:", perfile.Path, " cursize: ", cursize, "perser :", perfile.Size)
		}
	}
}

//MutiFileMonitor is
func MutiFileMonitor(pathList []string) {

	for _, path := range pathList {
		go FileListener(path)
	}

}

// FileListener 检测version
func FileListener(path string) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 2)

	// Set up a watchpoint listening for inotify-specific events within a
	// current working directory. Dispatch each InMovedFrom and InMovedTo
	// events separately to c.
	if err := notify.Watch(path, c, notify.Create, notify.InMovedFrom, notify.InModify, notify.InMovedTo, notify.Remove); err != nil {
		log.Error(err)
	}
	defer notify.Stop(c)

	// Inotify reports move filesystem action by sending two events tied with
	// unique cookie value (uint32): one of the events is of InMovedFrom type
	// carrying move source path, while the second one is of InMoveTo type
	// carrying move destination path.
	moves := make(map[uint32]struct {
		From string
		To   string
	})
	var flag = 0
	// Wait for moves.
	for ei := range c {
		cookie := ei.Sys().(*unix.InotifyEvent).Cookie

		info := moves[cookie]
		switch ei.Event() {
		case notify.InMovedFrom:
			info.From = ei.Path()
			//log.Warn("InMovedFrom", ei.Path())
		case notify.InMovedTo:
			info.To = ei.Path()
			//log.Warn("InMovedTo", ei.Path())
		case notify.Create:
			log.Info("create", ei.Path())
			if ei.Path() == "/edge/synctools.zip" {
				time.Sleep(300 * time.Microsecond)
				log.Info("has detecd the synctools package,")
				data, err := ioutil.ReadFile("/edge/synctools.zip")
				if err != nil {
					log.Error(err)
				}
				utils.Unzip(data, "/edge/mnt")
				exec.Command("has unzip the synctools","/edge/jxcore/bin/").Run()
				log.Info("If you want to use synctools ,please shutdown current process ,and modify the settings.yaml git")
			}

		case notify.Remove:
			//log.Warn("remove", ei.Path())
		case notify.InModify:

			if ei.Path() == "/edge/monitor/telegraf/bin/telegraf.cfg" || ei.Path() == "/edge/monitor/device-statsite/conf/influxdb.ini" {

				if flag == 0 {
					log.Info("Generate configuration ...")
				}
				flag += 1
				if flag == 6 {
					flag = 0
				}
			}

			if ei.Path()[len(ei.Path())-7:] == "version" {
				version.ChangLog(version.PraseVersionFile())
			}
			if ei.Path() == "/edge/synctools.zip" {
				time.Sleep(300 * time.Microsecond)
				log.Info("has detecd the synctools package modify,")
				data, err := ioutil.ReadFile("/edge/synctools.zip")
				if err != nil {
					log.Error(err)
				}
				utils.Unzip(data, "/edge/mnt")
				exec.Command("has unzip the synctools","/edge/jxcore/bin/").Run()
				log.Info("has unzip the synctools,")
				log.Info("If you want to use synctools ,please shutdown current process ,and modify the settings.yaml git")
			}

		}
		moves[cookie] = info

		if cookie != 0 && info.From != "" && info.To != "" {
			log.Warn("File:", info.From, "was renamed to", info.To)
			delete(moves, cookie)
		}
	}
}
