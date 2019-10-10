package updatemanage

import "sync"

const (
    FINISHED UpgradeStatus = iota
    UPDATING UpgradeStatus = 10
)
const (
    EDGEVERSIONFILE = "/edge/VERSION"
    TARGETVERSION   = "/etc/edgetarget"
    UPLOADURL       = "http://10.55.2.207:10111/api/v1/worker_version"
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
