package manager

import (
	"fmt"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"sync"
	"testing"
)

func TestInsert(t *testing.T) {
	aLog := logs.NewOplog(types.NETWORKE, "注册到master")
	bLog := logs.NewOplog(types.NETWORKE, "hello")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			Insert(aLog)
		}
		wg.Done()
	}(wg)
	for i := 0; i < 10; i++ {
		Insert(bLog)
	}
	wg.Wait()
	logs, err := ListAll()
	if err != nil {
		t.Error(err)
	}
	for _, log := range logs {
		fmt.Println(string(log.Marshal()))
	}

}
