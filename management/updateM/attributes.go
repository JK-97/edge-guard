package updateM

import "sync"

const (
    FINISHED     UpgradeStatus = iota
    UPDATING                   = 10
    UPDATESOURCE               = 100
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
