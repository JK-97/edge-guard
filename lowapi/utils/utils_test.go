package utils

import (
	"io/ioutil"
    log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"testing"
)



func TestUnzip(t *testing.T) {
	data,err:=ioutil.ReadFile("/home/marshen/synctools.zip")
	if err != nil {
		log.Error(err)
	}
	Unzip(data,"/home/marshen/new")
}
