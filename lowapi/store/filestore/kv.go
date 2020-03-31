package filestore

import (
	"github.com/JK-97/edge-guard/lowapi/store"
	"io/ioutil"
	"os"
	"path"
)

type kv struct {
	baseDir string
}

var KV store.KV = &kv{baseDir: "/edge/store"}

func (kv *kv) Get(key string) ([]byte, error) {
	data, err := ioutil.ReadFile(kv.getPath(key))
	return data, err
}

func (kv *kv) GetDefault(key string, value []byte) ([]byte, error) {
	data, err := kv.Get(key)
	if err != nil {
		if os.IsNotExist(err) {
			return value, nil
		} else {
			return nil, err
		}
	}
	return data, nil
}

var inited bool

func (kv *kv) Set(key string, value []byte) error {
	keyPath := kv.getPath(key)
	if !inited {
		if err := os.MkdirAll(path.Dir(keyPath), 0755); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(keyPath, value, 0755)
}

func (kv *kv) getPath(key string) string {
	return path.Join(kv.baseDir, key)
}
