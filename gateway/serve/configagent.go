package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"jxcore/gateway/option"
	log "jxcore/lowapi/logger"
)

// ConfigWatcher 监听配置变化
type ConfigWatcher interface {
	WatchKey(key string, timeout time.Duration)
	UnwatchKey(key string)
	ResponseChan() <-chan watchKeyResponse
	Stop()
}

// ConfigClient 读取配置
type ConfigClient interface {
	GetConfig(string, time.Duration, bool) (*watchKeyResponse, error)
	PutConfig(string, interface{}) error
}

type httpConfigClient struct {
	Client http.Client
}

// ConfigAgentHandler 转发 ConfigAgent 相关请求
type ConfigAgentHandler struct {
	option.ConfigAgentConfig
	DialContext  func(context.Context, string, string) (net.Conn, error)
	Upgrader     websocket.Upgrader
	Client       ConfigClient // *http.Client
	ReverseProxy http.Handler
}

type keysWatcher struct {
	keys      map[string]time.Duration
	Client    ConfigClient // *http.Client
	stoped    bool
	responseC chan watchKeyResponse
	mu        *sync.Mutex
}

// NewConfigAgentHandler 获取新的 ConfigAgentHandler 实例
func NewConfigAgentHandler(config option.ConfigAgentConfig) *ConfigAgentHandler {
	handler := new(ConfigAgentHandler)
	u, err := url.Parse(config.Host)
	if err != nil {
		panic(err)
	}
	handler.ConfigAgentConfig = config
	handler.ReverseProxy = httputil.NewSingleHostReverseProxy(u)

	var d net.Dialer

	handler.Client = &httpConfigClient{
		Client: http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return d.DialContext(ctx, "tcp", u.Host)
				},
			},
		},
	}

	return handler
}

// 对 Key 的操作
const (
	ActionGet     = 0 // 获取设置
	ActionPut     = 1
	ActionWatch   = 3
	ActionUnwatch = 4
)

type getKeyRequest struct {
	Action  int           `json:"action"`
	Key     string        `json:"key"`
	JSON    interface{}   `json:"json"`
	Timeout time.Duration `json:"timeout,omitempty"`
}

type watchKeyResponse struct {
	Key string `json:"key"`

	configAgentData
	// JSON interface{} `json:"json"`
}

type configAgentData struct {
	Index int    `json:"index"`
	Value string `json:"Value"`
}

type configAgentResponse struct {
	Data        configAgentData `json:"data"`
	Description string          `json:"desc,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newKeysWatcher(client ConfigClient) *keysWatcher {
	return &keysWatcher{
		keys:      make(map[string]time.Duration),
		Client:    client,
		stoped:    false,
		responseC: make(chan watchKeyResponse, 3),
		mu:        new(sync.Mutex),
	}
}

func (w *httpConfigClient) GetConfig(key string, timeout time.Duration, watch bool) (reply *watchKeyResponse, err error) {

	escapedKey := url.PathEscape(key)
	var path string
	if watch {
		path = "/api/v1/watch/config/" + escapedKey
	} else {
		path = "/api/v1/config/" + escapedKey
	}

	u := url.URL{
		Scheme:   "http",
		Host:     "edgegw.localhost",
		Path:     path,
		RawQuery: fmt.Sprintf("timeout=%d", int(float32(timeout)/float32(time.Second))),
	}
	client := w.Client

	resp, err := client.Do(&http.Request{
		Method: http.MethodGet,
		URL:    &u,
	})
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		if resp.Body != nil {
			if buf, err := ioutil.ReadAll(resp.Body); err == nil {
				log.Debug(string(buf))
				resp.Body.Close()
			}
		}
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		var buf []byte
		buf, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		config := configAgentResponse{}

		if err = json.Unmarshal(buf, &config); err == nil {
			reply = &watchKeyResponse{
				Key:             key,
				configAgentData: config.Data,
				// JSON: config.Data,
			}
		}
	}

	return
}

func (w *httpConfigClient) PutConfig(key string, v interface{}) error {

	escapedKey := url.PathEscape(key)
	var path string
	path = "/api/v1/config/" + escapedKey

	u := url.URL{
		Scheme: "http",
		Path:   path,
	}
	client := w.Client

	buff, err := json.Marshal(v)
	if err != nil {
		return err
	}
	resp, err := client.Do(&http.Request{
		Method: http.MethodPut,
		URL:    &u,
		Body:   ioutil.NopCloser(bytes.NewBuffer(buff)),
	})
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		var buf []byte
		buf, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(buf))
	}

	return nil
}

func (w *keysWatcher) shouldStop(key string) (timeout time.Duration, stop bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	stop = w.stoped
	if stop {
		return
	}
	timeout, ok := w.keys[key]
	stop = !ok
	return
}

func (w *keysWatcher) watchKeybg(key string) {

	for {
		timeout, stop := w.shouldStop(key)
		if stop {
			break
		}
		if reply, err := w.Client.GetConfig(key, timeout, true); err == nil {
			w.responseC <- *reply
		}
	}
}

func (w *keysWatcher) WatchKey(key string, timeout time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.keys[key]; !ok {
		w.keys[key] = timeout
		go w.watchKeybg(key)
	} else {
		w.keys[key] = timeout
	}
}

func (w *keysWatcher) UnwatchKey(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.keys, key)
}

func (w *keysWatcher) ResponseChan() <-chan watchKeyResponse {
	return w.responseC
}

func (w *keysWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.stoped = true
	close(w.responseC)
}

func handleRequest(ch <-chan *getKeyRequest, configWatcher ConfigWatcher, c *websocket.Conn, client ConfigClient) {
	c1 := configWatcher.ResponseChan()
	for {
		select {
		case req, ok := <-ch:
			if !ok {
				configWatcher.Stop()
				return
			}
			switch req.Action {
			case ActionGet:
				resp, err := client.GetConfig(req.Key, req.Timeout, false)
				if err != nil {
					c.WriteJSON(errorResponse{err.Error()})
				} else {
					c.WriteJSON(resp)
				}
			case ActionPut:
				if req.JSON != nil {
					client.PutConfig(req.Key, req.JSON)
				} else {
					c.WriteJSON(errorResponse{"missing argument 'json'"})
				}
			case ActionWatch:
				configWatcher.WatchKey(req.Key, req.Timeout)
			case ActionUnwatch:
				configWatcher.UnwatchKey(req.Key)
			default:
				c.WriteJSON(errorResponse{"unknown action"})
			}
		case resp := <-c1:
			c.WriteJSON(resp)
		}
	}
}

// ServeWebsocket 处理 Websocket 长连接请求
func (h *ConfigAgentHandler) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	ch := make(chan *getKeyRequest, 32)

	c.SetCloseHandler(func(code int, text string) error {
		close(ch)
		message := websocket.FormatCloseMessage(code, "")
		c.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	go handleRequest(ch, newKeysWatcher(h.Client), c, h.Client)

	req := getKeyRequest{}

	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}

		switch mt {
		case websocket.TextMessage:
			err := json.Unmarshal(message, &req)
			if err != nil {
				c.WriteJSON(errorResponse{err.Error()})
			} else if req.Key == "" {
				c.WriteJSON(errorResponse{"missing argument 'key'"})
			} else {
				if req.Timeout == 0 {
					req.Timeout = h.Timeout
				}

				ch <- &req
			}
		}
	}
}

func (h *ConfigAgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ReverseProxy.ServeHTTP(w, r)
}
