package utils

import log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"

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