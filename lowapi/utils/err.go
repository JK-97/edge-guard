package utils

import "jxcore/log"

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