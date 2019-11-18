package network

import (
	"testing"
)

func TestFindBestIFace(t *testing.T) {
	// err := switchIFace(findBestIFace())
	// log.Info(err)
	InitIFace()
	MaintainBestIFace()
	// route, _ := getGWRoute("eth0")
	// log.Info(route)
}
