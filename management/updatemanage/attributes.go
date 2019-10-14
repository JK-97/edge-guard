package updatemanage

import "sync"

const (
    FINISHED UpgradeStatus = iota
    UPDATING UpgradeStatus = 10
)
const (
    EDGEVERSIONFILE = "/edge/VERSION"
    TARGETVERSION   = "/etc/edgetarget"
    UPLOADDOMAIN    = "port30111.version-control.ffffffffffffffffffffffff.master.iotedge"
    UPLOADPATH      = "/api/v1/worker_version"
    SOURCEHOST      = "10.53.1.220"
)

func (p UpgradeStatus) String() string {
    switch p {
    case FINISHED:
        return "finished"
    case UPDATING:
        return "updating"
    default:
        return "unknow"
    }

}

type UpgradeStatus int
type VersionInfo map[string]string

var process *UpgradeProcess
var lock *sync.Mutex = &sync.Mutex{}
