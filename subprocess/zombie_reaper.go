// +build !windows

package subprocess

import (
	reaper "github.com/ochinchina/go-reaper"
)

func ReapZombie() {
	go reaper.Reap()
}
