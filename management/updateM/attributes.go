package updateM

import "sync"

const (
    FINISHED     UpgradStatus = iota
    UPDATING                  = 10
    UPDATESOURCE              = 100
)

func (p UpgradStatus) String() string {
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

type UpgradStatus int
type Versioninfo map[string]string



var process *UpgradeProcess
var lock *sync.Mutex = &sync.Mutex{}