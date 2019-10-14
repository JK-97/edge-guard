package updatemanage

import "sync"

const (
    FINISHED     UpgradeStatus = iota
    UPDATING                   = 10
    UPDATESOURCE               = 100
)
const (
    EDGEVERSIONFILE = "/edge/VERSION"
    TARGETVERSION="/etc/edgetarget"
    UPLOADURL = "http://10.55.2.207:10111/api/v1/worker_version"
    UPLOADPATH ="port30111.version-control.ffffffffffffffffffffffff.master.iotedge"
)

func (p UpgradeStatus) String() string {
    switch p {
    case FINISHED:
        return "FINISH"
    case UPDATESOURCE:
        return "UPDATESOURCE"
    case UPDATING:
        return "UPDATING"
    default:
        return "UNKNOWN"
    }

}

type UpgradeStatus int
type VersionInfo map[string]string

var process *UpgradeProcess
var lock *sync.Mutex = &sync.Mutex{}
