package serve

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	// "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"jxcore/gateway/log"
	"jxcore/lowapi/logger"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/patrickmn/go-cache"
)

var CamerApiPath string = "http://localhost:48082/api/v1/device/%s/command/%s"
var c = cache.New(5*time.Minute, 5*time.Minute)

const (
	// StatusCleaning cleaning
	StatusCleaning SwitchModelStatus = "cleaning"
	// StatusLoading loading
	StatusLoading SwitchModelStatus = "loading"
	// StatusPreheating preheating
	StatusPreheating SwitchModelStatus = "preheating"
	// StatusLoaded loaded
	StatusLoaded SwitchModelStatus = "loaded"
	// StatusFailed failed
	StatusFailed SwitchModelStatus = "failed"

	GrpcServerAddress = "127.0.0.1:50051"
)
const (
	// ResultSucceed 操作成功
	ResultSucceed = "succ"

	// ResultFailed 操作失败
	ResultFailed = "fail"
)

// AiServingHandler 处理 AI Serving 的服务调用
type AiServingHandler struct {
	ServingAddr *url.URL
}

type aiDetectBody struct {
	Path string `json:"path"`
}

type aiModel struct {
	Model  string `json:"model"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

// SwitchModelStatus 切换模型状态
type SwitchModelStatus string

type aiSwitchRequest struct {
	// Model specify the model that want to switch or load
	Model string `json:"model"`
	// Mode <"frozen", "unfrozen">, specify the model is a frozen model or unfrozen model
	Mode string `json:"mode"`
	// Preheat specify whether to preheat the session
	Preheat bool `json:"preheat"`
}

type inferenceLocalRequest struct {
	CamerID string `json:"camer_id"`
	Model   string `json:"model"`
	Version string `json:"version"`
	Save    bool   `json:"save"`
}

type device struct {
	Origin string `json:"origin"`
	Device string `json:"device"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type getCapturePathReponce struct {
	Device   string   `json:"device"`
	Origin   string   `json:"origin"`
	Readings []device `json:"readings"`
}

type inferenceRemoteRequset struct {
	Model   string `json:"model"`
	Version string `json:"version"`
	Base64  string `json:"base64"`
	Save    string `json:"save"`
}
type createAIBackend struct {
	Model   string `json:"model"`
	Version string `json:"version"`
}
type aiModelReply struct {
	Result string `json:"result"`
	// Model current serving model
	Model string `json:"model"`
	// Status indicate current status of model switching
	Status SwitchModelStatus `json:"status"`
	// Error error message when failed to load a model
	Error string `json:"error"`
}

type aiSwitchModelReply struct {
	// Status indicate current status of model switching
	Status string `json:"status"`
}

// NewAiServingHandler 创建新的 AiServingHandler 实例
func NewAiServingHandler(u string) *AiServingHandler {
	pURL, err := url.Parse(u)
	if err != nil {
		log.Errorln(err)
	}

	return &AiServingHandler{pURL}
}

// handleDetect 调用 AI 检测
func (h *AiServingHandler) handleDetect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
	// TODO: 调用 AI 检测
	contentLength := r.Header.Get(http.CanonicalHeaderKey("Content-Length"))

	if contentLength == "" || contentLength == "0" {
		ErrorWithCode(w, http.StatusBadRequest)
		return
	}

	body := r.Body
	defer body.Close()

	b := aiDetectBody{}
	buff, err := ioutil.ReadAll(body)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buff, &b); err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if b.Path == "" {
		Error(w, "Missing Argument 'path'", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(b.Path); os.IsNotExist(err) {
		Error(w, "File Not Exist", http.StatusNotFound)
		return
	}

	logger.Printf("AI Detect: [%s]\n", b.Path)

	r.Body = ioutil.NopCloser(bytes.NewReader(buff))
	_url := h.ServingAddr
	proxy := httputil.NewSingleHostReverseProxy(_url)
	proxy.ServeHTTP(w, r)
}

// getModels 获取模型列表
func (h *AiServingHandler) getModels(w http.ResponseWriter, r *http.Request) {

}

// switchModel 切换模型
func (h *AiServingHandler) switchModel(w http.ResponseWriter, r *http.Request) {

}

/*
使用map[string][]string
key = 模型名称 val = 运行相同模型的后台bids

request 进来 ， 先在cache 中获取当前运行对应模型的后台中 选一个 ，循环尝试，成功直接退出，失败删除缓存
cache 中都没成功 获取新的缓存状态
*/

func responceJson(w http.ResponseWriter, obj interface{}, stautsCode int) {
	w.WriteHeader(stautsCode)
	data, err := json.Marshal(obj)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write(data)
}

func unmarshalRequest(r *http.Request, httpRequest interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &httpRequest)
	if err != nil {
		return err
	}
	return nil
}

var config = consulapi.DefaultConfig()
var consulClient *consulapi.Client

func init() {
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Error(err)
	}
	consulClient = client
}

type registry struct {
	AIName       string `json:"ai_name"`
	HeartbeatURL string `json:"heartbeat_url"`
	ServiceURL   string `json:"service_url"`
}
type antiRegistry struct {
	AIName string `json:"ai_name"`
}
type registerResp struct {
	Result string `json:"result"`
}

func (h *AiServingHandler) deRegisterAIHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		ErrorWithCode(w, 400)
		return
	}

	req := &antiRegistry{}
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Error(err)
		ErrorWithCode(w, 400)
		return
	}
	err = consulClient.Agent().ServiceDeregister(strings.Join([]string{req.AIName, req.AIName}, "."))
	if err != nil {

		ErrorWithCode(w, 400)
		return
	}
	resp := &registerResp{
		Result: "success",
	}
	respData, err := json.Marshal(resp)
	w.Write(respData)

}

func (h *AiServingHandler) registerAIHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}

	req := &registry{}
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Error(err)
		ErrorWithCode(w, 400)
		return
	}

	aiService := "ai_service"

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = strings.Join([]string{req.AIName, req.AIName}, ".")
	registration.Tags = []string{aiService}
	registration.Name = req.AIName
	registration.Kind = consulapi.ServiceKind(aiService)
	registration.Address = req.ServiceURL
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           req.HeartbeatURL,
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		ErrorWithCode(w, 400)
		log.Error("register server error : ", err)
		return
	}
	resp := &registerResp{
		Result: "success",
	}
	respData, err := json.Marshal(resp)
	w.Write(respData)

}

func (h *AiServingHandler) handleDetectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
	aiName := r.Header.Get("ai_name")

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	catalog := client.Catalog()
	// 使用缓存，缓存超过maxage 会再去获取一次
	services, _, err := catalog.Service(aiName, "ai_service", &consulapi.QueryOptions{
		UseCache: true,
		MaxAge:   3 * time.Hour,
	})

	dstServiceURL := ""
	for _, service := range services {
		if service.ServiceName == aiName {
			dstServiceURL = service.Address
		}
	}

	if dstServiceURL == "" {
		ErrorNotFound(w)
	}

	proxyURL, err := url.Parse(dstServiceURL)
	if err != nil {
		ErrorNotFound(w)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}

func (h *AiServingHandler) handleLocalDetectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
	aiName := r.Header.Get("ai_name")
	imagePath := r.Header.Get("image_path")

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	imageData, err := ioutil.ReadFile(imagePath)
	r.Body = ioutil.NopCloser(bytes.NewReader(imageData))

	catalog := client.Catalog()
	// 使用缓存，缓存超过maxage 会再去获取一次
	services, _, err := catalog.Service(aiName, "ai_service", &consulapi.QueryOptions{
		UseCache: true,
		MaxAge:   3 * time.Hour,
	})
	dstServiceURL := ""
	for _, service := range services {
		if service.ServiceName == aiName {
			dstServiceURL = service.Address
		}
	}

	if dstServiceURL == "" {
		ErrorNotFound(w)
	}

	proxyURL, err := url.Parse(dstServiceURL)
	if err != nil {
		ErrorNotFound(w)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
func (h *AiServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// logger.Printf("In:\t%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	path := r.URL.Path

	switch path {
	case "/v1/detect":
		h.handleDetect(w, r)
		return
	case "v1/switch":
		// TODO: 切换模型
		switch r.Method {
		case http.MethodGet:
		case http.MethodPost:
		default:
			ErrorWithCode(w, http.StatusMethodNotAllowed)
		}
	case "/v1/register":
		h.registerAIHandler(w, r)
	case "/v1/remotedetect":
		h.handleDetectHTTP(w, r)
		return
	case "/v1/localdetect":
		h.handleLocalDetectHTTP(w, r)
		return
	case "/v1/deregister":
		h.deRegisterAIHandler(w, r)
		return
	}
	_url := h.ServingAddr
	proxy := httputil.NewSingleHostReverseProxy(_url)

	// logger.Printf("Out:\t%v %s %s\n", _url, r.Method, r.URL)
	proxy.ServeHTTP(w, r)
}
