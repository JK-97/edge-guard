package utils

import log "jxcore/lowapi/logger"

func CheckErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

func CheckWarn(err error) {
	if err != nil {
		log.Error(err)
	}
}
