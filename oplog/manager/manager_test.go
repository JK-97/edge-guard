package manager

import (
	"encoding/json"
	"fmt"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"sync"
	"testing"
)

func TestInsert(t *testing.T) {
	// err := os.RemoveAll(defaultLogPath)
	// if err != nil {
	// 	t.Error(err)
	// }
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
	data, err := json.Marshal(logs)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

}
