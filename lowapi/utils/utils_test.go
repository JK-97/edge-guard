package utils

import (
	"io/ioutil"
    log "jxcore/go-utils/logger"
	"testing"
)



func TestUnzip(t *testing.T) {
	data,err:=ioutil.ReadFile("/home/marshen/synctools.zip")
	if err != nil {
		log.Error(err)
	}
	Unzip(data,"/home/marshen/new")
}
