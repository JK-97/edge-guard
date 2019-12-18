package serve

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	// "log"
	"net/http"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	log "jxcore/lowapi/logger"
	"jxcore/gateway/store"
)

// DynamicService 动态服务
type DynamicService struct {
	Name        string   `json:"name"`
	HostPort    string   `json:"host_port"`
	Version     string   `json:"version"`
	HealthCheck string   `json:"health_check"`
	URLPatterns []string `json:"url_patterns"`
}

// DynamicServiceContainer 动态服务容器
type DynamicServiceContainer struct {
	Services map[string][]*DynamicService

	// OnChange 服务数量发生变化时调用
	OnChange func()

	mu sync.Locker

	Prefix string
}

// DynamicServiceHandler 动态服务处理
type DynamicServiceHandler struct {
	Jar *DynamicServiceContainer
}

// NewDynamicServiceContainer 动态服务容器
// var DynamicServiceJar = DynamicServiceContainer{}
func NewDynamicServiceContainer() *DynamicServiceContainer {
	return &DynamicServiceContainer{
		Services: make(map[string][]*DynamicService, 10),
		mu:       new(sync.Mutex),
		Prefix:   "ds-",
	}
}

// Add 服务注册
func (d *DynamicServiceContainer) Add(ds *DynamicService) {
	d.mu.Lock()
	defer d.mu.Unlock()
	name := ds.Name
	list := d.Services[name]

	if list == nil {
		d.Services[name] = []*DynamicService{ds}
	} else {
		for i, it := range list {
			if it.HostPort == ds.HostPort {
				// 新注册的服务放在最前面
				if i > 0 {
					list[0], list[i] = list[i], list[0]
					d.Services[name] = list
				}
				if d.OnChange != nil {
					d.OnChange()
				}
				return
			}
		}
		// 添加服务
		length := len(list)
		list = append(list, ds)
		list[0], list[length] = list[length], list[0]
		d.Services[name] = list
	}
	if d.OnChange != nil {
		d.OnChange()
	}
}

// RemoveByHostPort 服务反注册
func (d *DynamicServiceContainer) RemoveByHostPort(s string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for name, arr := range d.Services {
		for i, it := range arr {
			if it.HostPort == s {
				upperBound := len(arr)
				if i == 0 {
					d.Services[name] = arr[1:upperBound]
				} else if i == upperBound {
					d.Services[name] = arr[0 : upperBound-1]
				} else {
					d.Services[name] = append(arr[0:i], arr[i+1:upperBound]...)
				}

				if d.OnChange != nil {
					d.OnChange()
				}
				return
			}
		}
	}
}

// Remove 服务反注册
func (d *DynamicServiceContainer) Remove(ds *DynamicService) {
	d.RemoveByHostPort(ds.HostPort)
}

func (d *DynamicServiceContainer) Write(w io.Writer) {
	enc := toml.NewEncoder(w)

	enc.Encode(*d)
}

// Store 保存现有配置
func (d *DynamicServiceContainer) Store(s store.Store) {
	for name, srv := range d.Services {
		var lastKey []byte
		prefix := fmt.Sprintf(d.Prefix + name)
		for i, service := range srv {
			k := []byte(fmt.Sprintf("%s.%d", prefix, i))
			b, err := json.Marshal(service)
			if err != nil {
				continue
			}
			s.Put(k, b)

			lastKey = k
		}
		if lastKey != nil {
			k := string(lastKey)
			iter := s.NewIterator(prefix)
			for iter.Next() {

				if string(iter.Key()) <= k {
					continue
				}
				s.Delete(iter.Key())
			}
			iter.Release()
		}
	}
}

// Recovery 从存储中恢复配置
func (d *DynamicServiceContainer) Recovery(s store.Store) {
	iter := s.NewIterator(d.Prefix)
	for iter.Next() {
		// json.
		// iter.Value()
		srv := DynamicService{}
		if err := json.Unmarshal(iter.Value(), &srv); err != nil {
			log.Println(err)
			continue
		}
		d.Add(&srv)
	}
	iter.Release()
}

// Search 服务查找
func (d *DynamicServiceContainer) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	version := query.Get("version")

	var services []*DynamicService

	if name != "" {
		services = d.Services[name]
		if services == nil {
			services = make([]*DynamicService, 0)
		}
	} else {
		services = make([]*DynamicService, 0)
		for _, arr := range d.Services {
			services = append(services, arr...)
		}
	}

	if version != "" {
		length := len(services)
		left := 0
		right := length - 1

		for {
			if left >= right {
				break
			}
			if services[left].Version != version {
				services[left], services[right] = services[right], services[left]
				right--
			}
			left++
		}
		if services[left].Version != version {
			left--
		}
		services = services[0:left]
	}
	data := make(map[string]interface{})
	data["results"] = services
	data["total"] = len(services)
	WriteData(w, &data)
}

// Register 服务注册
func (d *DynamicServiceContainer) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		ErrorWithCode(w, http.StatusBadRequest)
	}

	services := []DynamicService{}
	err = json.Unmarshal(body, &services)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, srv := range services {
		d.Add(&srv)
	}

	WriteSucess(w)
}

// Unregister 服务反注册
func (d *DynamicServiceContainer) Unregister(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		ErrorWithCode(w, http.StatusBadRequest)
	}

	service := DynamicService{}
	err = json.Unmarshal(body, &service)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
	}

	d.Remove(&service)

	WriteSucess(w)
}

func (d *DynamicServiceContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.Search(w, r)
	case http.MethodPost:
		d.Register(w, r)
	case http.MethodDelete:
		d.Unregister(w, r)
	default:
		ErrorWithCode(w, http.StatusMethodNotAllowed)
	}
}

// DynamicHandler 处理动态服务的转发
func (d *DynamicServiceContainer) DynamicHandler() *DynamicServiceHandler {
	return &DynamicServiceHandler{d}
}

func (d *DynamicServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var name = r.Header.Get("X-Internal-Service")

	if name == "" {
		path := r.URL.Path

		name = strings.Split(path, "/")[0]

		if name == "" {
			ErrorNotFound(w)
			return
		}
	}

}
