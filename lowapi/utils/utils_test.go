package utils

import (
	log "github.com/JK-97/edge-guard/lowapi/logger"
	"io/ioutil"
	"testing"
)

func TestUnzip(t *testing.T) {
	data, err := ioutil.ReadFile("/home/marshen/synctools.zip")
	if err != nil {
		log.Error(err)
	}
	Unzip(data, "/home/marshen/new")
}
