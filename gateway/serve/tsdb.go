package serve

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/JK-97/edge-guard/gateway/option"
	"github.com/JK-97/edge-guard/gateway/store"
)

// TSDBHandler 处理 时序数据库相关请求
type TSDBHandler struct {
	option.TimeSeriesDBConfig
	Store store.Store
	mu    sync.Locker
}

// NewTSDBHandler 由给定的配置获取 TSDBHandler 实例
func NewTSDBHandler(config option.TimeSeriesDBConfig) *TSDBHandler {
	h := &TSDBHandler{
		TimeSeriesDBConfig: config,
		mu:                 new(sync.Mutex),
	}
	if h.StorePrefix == "" {
		h.StorePrefix = "/tsdb"
	}
	return h
}

// SetStore 设置保存配置用的 Store 实例
func (h *TSDBHandler) SetStore(s store.Store) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Store = s
}

// AddDatabase 添加要同步的 DB
func (h *TSDBHandler) AddDatabase(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Store == nil {
		return
	}
	key := h.StorePrefix + name
	h.Store.Put([]byte(key), nil)
}

// RemoveDatabase 从同步的 DB 中移除指定的 DB
func (h *TSDBHandler) RemoveDatabase(name string) {

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Store == nil {
		return
	}
	key := h.StorePrefix + name
	h.Store.Delete([]byte(key))
}

// ListDatabase 获取需要同步的 Database 列表
func (h *TSDBHandler) ListDatabase() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]string, 0)
	if h.Store == nil {
		return result
	}

	iter := h.Store.NewIterator(h.StorePrefix)
	prefixLen := len([]byte(h.StorePrefix))
	for iter.Next() {
		key := iter.Key()[prefixLen:]
		result = append(result, string(key))
	}

	return result
}

// CreateInfluxDBURI 获取 uri
func (h *TSDBHandler) CreateInfluxDBURI(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	localAddr := PickLocalAddr(r)
	if localAddr == "127.0.0.1" {
		localAddr = "172.17.0.1"
	}
	data["uri"] = net.JoinHostPort(localAddr, strconv.Itoa(h.InfluxDBPort))
	data["user"] = h.InfluxDBUser
	data["password"] = h.InfluxDBPassword

	WriteData(w, &data)
}

// CreateStatsdURI 获取 uri
func (h *TSDBHandler) CreateStatsdURI(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	localAddr := PickLocalAddr(r)
	if localAddr == "127.0.0.1" {
		localAddr = "172.17.0.1"
	}

	data["statsite"] = net.JoinHostPort(localAddr, strconv.Itoa(h.StatdPort))

	WriteData(w, &data)
}

// AddSyncDatabaseHTTP 添加同步数据库
func (h *TSDBHandler) AddSyncDatabaseHTTP(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorWithCode(w, http.StatusBadRequest)
		return
	}
	var nameStruct struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(buf, &nameStruct)
	if err != nil {
		ErrorWithCode(w, http.StatusBadRequest)
		return
	}

	h.AddDatabase(nameStruct.Name)
	WriteSucess(w)
}

// RemoveSyncDatabaseHTTP 移除同步数据库
func (h *TSDBHandler) RemoveSyncDatabaseHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		Error(w, "Missing Argument `name`", http.StatusBadRequest)
		return
	}
	h.RemoveDatabase(name)

	WriteSucess(w)
}

// ListSyncDatabaseHTTP 获取同步数据库列表
func (h *TSDBHandler) ListSyncDatabaseHTTP(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	data["result"] = h.ListDatabase()

	WriteData(w, &data)
}

func (h *TSDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/statsaddr":
		switch r.Method {
		case http.MethodPost:
			h.CreateStatsdURI(w, r)
		default:
			ErrorWithCode(w, http.StatusMethodNotAllowed)
		}
		return
	case "/influxaddr":
		switch r.Method {
		case http.MethodPost:
			h.CreateInfluxDBURI(w, r)
		default:
			ErrorWithCode(w, http.StatusMethodNotAllowed)
		}
		return
	case "/sync":
		switch r.Method {
		case http.MethodGet:
			h.ListSyncDatabaseHTTP(w, r)
		case http.MethodPost:
			h.AddSyncDatabaseHTTP(w, r)
		case http.MethodDelete:
			h.RemoveSyncDatabaseHTTP(w, r)
		default:
			ErrorWithCode(w, http.StatusMethodNotAllowed)
		}
		return
	}

	ErrorWithCode(w, http.StatusNotFound)
}
