package signal

import (
	"os"
)

func init() {
	StopSignals = []os.Signal{
		os.Interrupt,
		os.Kill,
	}
}
