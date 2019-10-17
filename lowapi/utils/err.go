package utils

import log "jxcore/go-utils/logger"

func CheckErr(err error)  {
    if err != nil {
        log.Error(err)
    }
}

func CheckWarn(err error)  {
    if err != nil {
        log.Error(err)
    }
}