package process

import (

	"testing"
)

func TestNewProcess(t *testing.T) {
	p := NewProcess("1", "/edge/monitor/telegraf/bin/telegraf")
	p.Start(true)
	a := 1
	for {
		a++
	}

}
