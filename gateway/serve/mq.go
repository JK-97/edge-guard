package serve

import (

	// "log"
	"container/list"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"jxcore/gateway/log"
	"jxcore/gateway/option"
	"jxcore/gateway/store"
)

const jsonContentType = "application/json"

// MesageQueueInternalHandler 处理 MQ 相关调用请求
type MesageQueueInternalHandler struct {
	// MessageQueueURI string
	// PorterSocket    bool
	option.MessageQueueConfig

	stroe store.Store

	porterNotifier PorterNotifier

	Upgrader *websocket.Upgrader

	mu      *sync.Mutex
	wsConns *list.List
}

// PorterNotifier 通知 Porter Topic 变化
type PorterNotifier interface {
	// GrantTopic 获取 topic 授权/取消 topic 授权
	GrantTopic(topic string, grant bool) error

	// ListTopics 枚举 Topic
	ListTopics() []string
}

// NewMesageQueueInternalHandler 创建 MesageQueueInternalHandler 实例
func NewMesageQueueInternalHandler(config option.MessageQueueConfig) *MesageQueueInternalHandler {
	if config.MessageQueueURI == "" {
		config.MessageQueueURI = "amqp://guest:guest@localhost:5672/%2F"
	}
	if config.MessageQueuePrefix == "" {
		config.MessageQueuePrefix = "/topicz/"
	}
	var porterNotifier PorterNotifier

	porterNotifier = &StorePorterNotifier{
		mu:     new(sync.Mutex),
		prefix: config.MessageQueuePrefix,
	}
	// if config.PorterSocket != "" {
	// 	porterNotifier = NewUnixPorterNotifier(config.PorterSocket)
	// }

	return &MesageQueueInternalHandler{
		MessageQueueConfig: config,
		porterNotifier:     porterNotifier,
		mu:                 new(sync.Mutex),
		Upgrader:           new(websocket.Upgrader),
		wsConns:            new(list.List),
	}
}

// StorePorterNotifier 基于存储的 Porter 通知器
type StorePorterNotifier struct {
	store  store.Store
	prefix string
	mu     *sync.Mutex
}

type grantTopicRequest struct {
	Topics []string `json:"topics"`
}

// const prefix = "/topicz/"

// GrantTopic 获取 topic 授权/取消 topic 授权
func (p *StorePorterNotifier) GrantTopic(topic string, grant bool) error {
	if p.store == nil {
		return nil
	}
	if !grant {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	prefix := p.prefix
	key := prefix + topic
	err := p.store.Put([]byte(key), nil)

	return err
}

// ListTopics 枚举 Topic
func (p *StorePorterNotifier) ListTopics() []string {
	result := make([]string, 0, 8)
	if p.store == nil {
		return result
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	prefix := p.prefix
	iter := p.store.NewIterator(prefix)
	prifxLen := len([]byte(prefix))
	for iter.Next() {
		topic := string(iter.Key()[prifxLen:])
		result = append(result, topic)
	}

	return result
}

// GrantTopic 授权访问 MQ 上的 Topic
func (h *MesageQueueInternalHandler) GrantTopic(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	topic := data["topic"]
	if topic == "" {
		topic = r.URL.Query().Get("topic")
	}
	log.Printf("GrantTopic: %s", topic)

	if h.porterNotifier != nil {
		h.porterNotifier.GrantTopic(topic, true)
	}
	WriteSucess(w)
	go h.sendTopic(topic, true)
}

// ProhibitTopic 取消对 MQ 上的 Topic 的授权
func (h *MesageQueueInternalHandler) ProhibitTopic(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	topic := data["topic"]
	if topic == "" {
		topic = r.URL.Query().Get("topic")
	}
	log.Printf("ProhibitTopic: %s", topic)
	if h.porterNotifier != nil {
		// 通知对应的服务，有服务取消对 Topic 的引用
		h.porterNotifier.GrantTopic(topic, false)
	}
	// go h.sendTopic(topic, false)
	WriteSucess(w)
}

// CreateMessageQueueURI 将
func (h *MesageQueueInternalHandler) CreateMessageQueueURI(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	uri := h.MessageQueueURI
	loc := PickLocalAddr(r)
	data["uri"] = strings.Replace(uri, "localhost", loc, 1)
	WriteData(w, &data)
}

// ListTopics TODO: 获取当前节点上所有的 Topic
func (h *MesageQueueInternalHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["topics"] = h.porterNotifier.ListTopics()
	WriteData(w, &data)
}

// addWsConn 添加长连接
func (h *MesageQueueInternalHandler) addWsConn(conn *websocket.Conn) *list.Element {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.wsConns.PushBack(conn)
}

// removeWsConn 移除长连接
func (h *MesageQueueInternalHandler) removeWsConn(ele *list.Element) {
	h.mu.Lock()
	defer h.mu.Unlock()

	l := h.wsConns
	l.Remove(ele)
}

type topicToken struct {
	Topic string `json:"topic"`
	Grant bool   `json:"grant"`
}

func (h *MesageQueueInternalHandler) sendTopic(topic string, grant bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	l := h.wsConns
	for e := l.Front(); e != nil; e = e.Next() {
		if c, ok := e.Value.(*websocket.Conn); ok {
			c.WriteJSON(topicToken{Topic: topic, Grant: grant})
		}
	}
}

// NotifyTopics 保持长连接，监听需同步的消息主题的变化 Websockets
func (h *MesageQueueInternalHandler) NotifyTopics(w http.ResponseWriter, r *http.Request) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	ele := h.addWsConn(c)

	c.SetCloseHandler(func(code int, text string) error {
		h.removeWsConn(ele)
		message := websocket.FormatCloseMessage(code, "")
		c.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	topics := h.porterNotifier.ListTopics()
	var token topicToken = topicToken{Grant: true}
	for _, topic := range topics {
		token.Topic = topic
		c.WriteJSON(&token)
	}
	c.WriteJSON(topics)

	for {
		c.ReadMessage()
		// mt, message, err := c.ReadMessage()
	}
}

// SetStore 设置保存配置用的 Store 实例
func (h *MesageQueueInternalHandler) SetStore(s store.Store) {
	h.stroe = s
	h.porterNotifier = &StorePorterNotifier{
		mu:     new(sync.Mutex),
		store:  s,
		prefix: h.MessageQueuePrefix,
	}
}

func (h *MesageQueueInternalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/grant") {
		switch r.Method {
		case "POST":
			h.GrantTopic(w, r)
		case "DELETE":
			h.ProhibitTopic(w, r)
		default:
			ErrorWithCode(w, http.StatusMethodNotAllowed)
		}
	} else if path == "/create" {
		h.CreateMessageQueueURI(w, r)
	} else if path == "/topicz" {
		h.ListTopics(w, r)
	} else if path == "/ws" {
		h.NotifyTopics(w, r)
	} else {
		ErrorNotFound(w)
	}
}
