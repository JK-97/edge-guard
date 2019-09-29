package utils

import (
	"io/ioutil"
	"jxcore/log"
	"testing"
)



func TestUnzip(t *testing.T) {
	data,err:=ioutil.ReadFile("/home/marshen/synctools.zip")
	if err != nil {
		log.Error(err)
	}
	Unzip(data,"/home/marshen/new")
}
