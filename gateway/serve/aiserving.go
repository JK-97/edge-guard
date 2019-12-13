package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	// "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"jxcore/gateway/log"
	pb "jxcore/gateway/trueno"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

var CamerApiPath string = "http://localhost:48082/api/v1/device/%s/command/%s"
var c = cache.New(5*time.Minute, 5*time.Minute)

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

const (
	// ResultSucceed 操作成功
	ResultSucceed = "succ"

	// ResultFailed 操作失败
	ResultFailed = "fail"
)

// SwitchModelStatus 切换模型状态
type SwitchModelStatus string

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
)

type aiSwitchRequest struct {
	// Model specify the model that want to switch or load
	Model string `json:"model"`
	// Mode <"frozen", "unfrozen">, specify the model is a frozen model or unfrozen model
	Mode string `json:"mode"`
	// Preheat specify whether to preheat the session
	Preheat bool `json:"preheat"`
}

type aiDetectRequest struct {
	CamerID string `json:"camer_id"`
	Model   string `json:"model"`
	Save    bool   `json:"save"`
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
		log.Fatalln(err)
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

	log.Printf("AI Detect: [%s]\n", b.Path)

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

func (h *AiServingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Printf("In:\t%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
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
	case "/v2/detect":
		h.aiDetection(w, r)
		return
	}
	_url := h.ServingAddr
	proxy := httputil.NewSingleHostReverseProxy(_url)

	// log.Printf("Out:\t%v %s %s\n", _url, r.Method, r.URL)
	proxy.ServeHTTP(w, r)
}

// 获取正在运行的模型
func grpcRunningBackend(conn *grpc.ClientConn, ctx context.Context) ([]*pb.RunningReply_Status, error) {
	client := pb.NewBackendClient(conn)
	resoponse, err := client.ListRunningBackends(ctx, &pb.PingRequest{Client: "client"})
	if err != nil {
		return nil, err
	}
	reply := resoponse.GetStatus()
	return reply, nil
}

// 本地检测
func grpcInferenceLocal(conn *grpc.ClientConn, InferRequest *pb.InferRequest, ctx context.Context) (*pb.ResultReply, error) {
	client := pb.NewInferenceClient(conn)
	reply, err := client.InferenceLocal(ctx, InferRequest)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// 更新缓存
func updateCache(resply []*pb.RunningReply_Status) {
	thisModelCache := []string{}
	for _, backend := range resply {
		if cache, ok := c.Get(backend.GetModel()); ok {
			if result, ok := cache.([]string); ok {
				thisModelCache = result
			}

		}
		thisModelCache = append(thisModelCache, backend.GetBid())
		if len(thisModelCache) == 0 {
			c.Delete(backend.GetModel())
			return
		}
		c.Set(backend.GetModel(), thisModelCache, 5*time.Minute)
	}
}

//删除模型对应的bid
func removeBidCache(model, bid string, cache []string) {
	newCache := []string{}
	for _, beackendId := range cache {
		if bid != beackendId {
			newCache = append(newCache, beackendId)
		}
	}
	c.Set(model, newCache, 5*time.Minute)
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

//通过camerid 获取 capture path
func getCapturePathByCamId(camId string) (string, error) {
	responce, err := http.Get(fmt.Sprintf(CamerApiPath, camId, camId))
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		return "", err
	}
	capturePathReponce := getCapturePathReponce{}
	err = json.Unmarshal(data, &capturePathReponce)
	if err != nil {
		return "", err
	}
	for _, device := range capturePathReponce.Readings {
		if device.Name == "capture_path" {
			valuemap := map[string]string{}
			err := json.Unmarshal([]byte(device.Value), &valuemap)
			if err != nil {
				return "", err
			}
			return valuemap["capture_path"], nil

		}
	}
	return "", errors.New("找不到摄像头")
}

//通过model名字获取对应bids
func getBackendByModel(conn *grpc.ClientConn, ctx context.Context, model string) ([]string, error) {
	result, ok := c.Get(model)
	if !ok {
		// 没有缓存先更新下缓存
		resply, err := grpcRunningBackend(conn, ctx)
		if err != nil {
			return nil, err
		}
		updateCache(resply)
		result, ok = c.Get(model)
		if ok {
			//再次检查，若没有可用的缓存，直接返回
			return nil, errors.New("未找到可用的 对应model 可用的backend")
		}
	}
	bidsResult, ok := result.([]string)
	if !ok {
		return nil, errors.New("bid 数据格式错误")
	}
	return bidsResult, nil

}

/*

使用map[string][]string
key = 模型名称 val = 运行相同模型的后台bids

request 进来 ， 先在cache 中获取当前运行对应模型的后台中 选一个 ，循环尝试，成功直接退出，失败删除缓存
cache 中都没成功 获取新的缓存状态

*/

func (h *AiServingHandler) aiDetection(w http.ResponseWriter, r *http.Request) {
	// TODO: grpc AI 本地识别
	//httprequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	httpRequest := aiDetectRequest{}
	err = json.Unmarshal(data, &httpRequest)
	if err != nil {
		return
	}
	//ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//rpc
	conn, err := grpc.Dial("tcp", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	//获取 cam capture  path
	capturePath, err := getCapturePathByCamId(httpRequest.CamerID)
	if err != nil {
		respoceJson(w, err, 400)
		return
	}

	//通过model 名字获取backend 的bid
	bidsResult, err := getBackendByModel(conn, ctx, httpRequest.Model)
	if err != nil {
		respoceJson(w, err, 400)
		return
	}

	//有缓存则尝试进行请求
	//从第一个开始尝试
	for _, bid := range bidsResult {
		reply, err := grpcInferenceLocal(conn, &pb.InferRequest{
			Bid:  bid,
			Uuid: uuid.New().String(),
			Path: capturePath,
			Type: "",
		}, ctx)
		if err != nil {
			continue
		}
		if reply.GetCode() != 0 {
			//如果失败，则删除这个bid缓存，尝试下一个bid 缓存
			removeBidCache(httpRequest.Model, bid, bidsResult)
			continue
		}
		respoceJson(w, reply, 400)
		return
	}

	respoceJson(w, []byte("not find useable backend in map"), 200)
	//都没成功 更新状态，直接返回

}

func respoceJson(w http.ResponseWriter, obj interface{}, stautsCode int) {
	w.WriteHeader(stautsCode)
	data, err := json.Marshal(obj)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write(data)
}
