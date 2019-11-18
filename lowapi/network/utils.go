package network

import (
	"hash/fnv"
	"jxcore/gateway/log"
)

var (
	// ifacePriorityHash 用来当iface route的metric, 取值200-400
	ifacePriorityHash = map[string]int{}
)

func hashStrToInt(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32()%200 + 200)
}

func initIFacePriorityHash() {
	occupied := map[int]bool{}
	for _, iFace := range ifacePriority {
		h := hashStrToInt(iFace)
		for occupied[h] {
			h += 1
		}
		ifacePriorityHash[iFace] = h
		occupied[h] = true
	}
	log.Debug("init ifacePriorityHash: %v", ifacePriorityHash)
}
