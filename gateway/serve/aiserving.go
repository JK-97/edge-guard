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

	"jxcore/gateway/dao"
	"jxcore/gateway/log"
	pb "jxcore/gateway/trueno"
	"jxcore/lowapi/logger"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
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

// 本地识别
func (h *AiServingHandler) aiLocalDetection(w http.ResponseWriter, r *http.Request) {
	//httprequest
	httpRequest := &inferenceLocalRequest{}
	err := unmarshalRequest(r, &httpRequest)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	//ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//rpc
	conn, err := grpc.Dial(GrpcServerAddress, grpc.WithInsecure())
	if err != nil {
		logger.Info(err.Error())
		return
	}
	defer conn.Close()

	//获取 cam capture  path
	capturePath := "/capture/"
	capturePath, err = getCapturePathByCamId(httpRequest.CamerID)
	if err != nil {
		responceJson(w, err, 400)
		return
	}

	//通过model 名字获取backend 的bid
	err = tryEveryBackend(conn, httpRequest.Model, httpRequest.Version, w, func(bid, detectUuid string) (*pb.ResultReply, error) {
		return grpcInferenceLocal(conn, &pb.InferRequest{
			Bid:  bid,
			Uuid: detectUuid,
			Path: capturePath,
			Type: "",
		}, ctx)
	}, ctx)
	if err != nil {
		responceJson(w, err.Error(), 200)
	}
}

//远程识别
func (h *AiServingHandler) aiRemoteDetection(w http.ResponseWriter, r *http.Request) {
	//httprequest
	httpRequest := &inferenceRemoteRequset{}
	err := unmarshalRequest(r, &httpRequest)
	if err != nil {
		responceJson(w, err.Error(), 400)
		return
	}
	//ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.Dial(GrpcServerAddress, grpc.WithInsecure())
	if err != nil {
		logger.Info(err.Error())
		return
	}
	defer conn.Close()
	err = tryEveryBackend(conn, httpRequest.Model, httpRequest.Version, w, func(bid, detectUuid string) (*pb.ResultReply, error) {
		return grpcInferenceRemote(conn, &pb.InferRequest{
			Bid:    bid,
			Uuid:   detectUuid,
			Base64: httpRequest.Base64,
			Type:   "",
		}, ctx)
	}, ctx)
	if err != nil {
		responceJson(w, err.Error(), 200)
	}

}

// 创建对应模型的后台
func (h *AiServingHandler) createAIBackend(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	httpRequest := inferenceLocalRequest{}
	err = json.Unmarshal(data, &httpRequest)
	if err != nil {
		responceJson(w, err.Error(), 400)
		return
	}
	//rpc
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		responceJson(w, err.Error(), 400)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := grpcCreateAndLoadModel(conn, ctx, httpRequest.Model, httpRequest.Version)
	if err != nil {
		responceJson(w, err.Error(), 400)
		return
	}

	logger.Info("创建模型后台:", httpRequest.Model, httpRequest.Version)
	logger.Info("创建后台bid:", msg)
	responceJson(w, "success", 200)
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
	return "", errors.New("can not find camera device")
}

func tryEveryBackend(conn *grpc.ClientConn, model string, version string, w http.ResponseWriter, inference func(bid, detectUuid string) (*pb.ResultReply, error), ctx context.Context) error {
	bidsResult, err := grpcGetBackendByModel(conn, ctx, model)
	if err != nil {
		_, _ = grpcCreateAndLoadModel(conn, ctx, model, version)
		return errors.New("自动创建model后台,请重试")
	}
	logger.Info("检索到的后台bid:", bidsResult)
	//有缓存则尝试进行请求
	//从第一个开始尝试
	detectUuid := uuid.New().String()
	for _, bid := range bidsResult {
		reply, err := inference(bid, detectUuid)
		if err != nil {
			continue
		}
		if reply.GetCode() != 0 {
			//如果失败，则删除这个bid缓存，尝试下一个bid 缓存
			removeBidCache(model, bid, bidsResult)
			continue
		}
		//通过uuid获取redis 数据
		redis, err := dao.NewRedisClient()
		if err != nil {
			return err
		}
		result, err := redis.Get(detectUuid).Result()
		if err != nil {
			return err
		}
		responceJson(w, result, 200)
		return nil
	}
	return errors.New("not find useable backend in map")
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
		logger.Info("更新缓存:", thisModelCache)
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

// 获取正在运行的模型
func grpcRunningBackend(conn *grpc.ClientConn, ctx context.Context) ([]*pb.RunningReply_Status, error) {
	client := pb.NewBackendClient(conn)
	resoponse, err := client.ListRunningBackends(ctx, &pb.PingRequest{Client: "client"})
	if err != nil {
		return nil, err
	}
	reply := resoponse.GetStatus()
	logger.Info("正在运行的模型后台连列表", reply)
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

// 远程检测
func grpcInferenceRemote(conn *grpc.ClientConn, InferRequest *pb.InferRequest, ctx context.Context) (*pb.ResultReply, error) {
	client := pb.NewInferenceClient(conn)
	reply, err := client.InferenceRemote(ctx, InferRequest)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// 列举存在的模型
func grpcListStoreModel(conn *grpc.ClientConn, ctx context.Context) ([]*pb.ModelInfo, error) {
	client := pb.NewModelClient(conn)
	resoponse, err := client.ListStoredModel(ctx, &pb.PingRequest{Client: "client"})
	if err != nil {
		return nil, err
	}
	reply := resoponse.GetList()
	logger.Info("model info", reply)
	return reply, nil
}

// 创建加载模型
func grpcCreateAndLoadModel(conn *grpc.ClientConn, ctx context.Context, model, version string) (string, error) {
	client := pb.NewInferenceClient(conn)
	resoponse, err := client.CreateAndLoadModel(ctx, &pb.LoadRequest{
		Bid:       "",
		Btype:     "tensorflow",
		Model:     model,
		Version:   version,
		Mode:      "frozen",
		Encrypted: 0,
	})
	if err != nil {
		return "", err
	}
	if resoponse.GetCode() != 0 {
		return "", errors.New("出现错误")
	}
	return resoponse.GetMsg(), nil
}

//通过model名字获取对应bids
func grpcGetBackendByModel(conn *grpc.ClientConn, ctx context.Context, model string) ([]string, error) {
	result, ok := c.Get(model)
	logger.Info("缓存结果：", ok, result)
	if !ok {
		// 没有缓存先更新下缓存
		resply, err := grpcRunningBackend(conn, ctx)
		if err != nil {
			return nil, err
		}
		updateCache(resply)
		result, ok = c.Get(model)
		if !ok {
			//再次检查，若没有可用的缓存，直接返回
			return nil, errors.New("can not find ai beckend")
		}
	}
	bidsResult, ok := result.([]string)
	if !ok {
		if len(bidsResult) == 0 {
			return nil, errors.New("can not find ai beckend ")
		}
		return nil, errors.New("bid data format err")
	}
	return bidsResult, nil

}

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

type registry struct {
	AIName       string `json:"ai_name"`
	HeartbeatURL string `json:"heartbeat_url"`
	ServiceURL   string `json:"service_url"`
}

// ----------------------------------------------------------------------------------------
// cache 存储的结构
var consulServiceCache = cache.New(5*time.Minute, 5*time.Minute)

func (h *AiServingHandler) registerAIHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}

	resp := &registry{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		log.Error(err)
		return
	}

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	uuid, err := uuid.NewUUID()
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = uuid.String()
	registration.Tags = []string{"ai_service"}
	registration.Name = resp.AIName
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           resp.HeartbeatURL,
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	registration.
		err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Error("register server error : ", err)
		return
	}
	w.Write([]byte("success"))

}

func (h *AiServingHandler) handleDetectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
	aiName := r.Header.Get("ai_name")
	raw, ok := consulServiceCache.Get(aiName)
	if !ok {

		ErrorNotFound(w)
	}
	serverURL, ok := raw.(string)
	if !ok {
		ErrorNotFound(w)
		return
	}
	proxyURL, err := url.Parse(serverURL)
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
	case "/v2/localdetect":
		h.aiLocalDetection(w, r)
	case "/v2/remotedetect":
		h.aiRemoteDetection(w, r)
	case "/v2/create":
		h.createAIBackend(w, r)

	case "/v3/register":
		h.registerAIHandler(w, r)

	case "/v3/remotedetect":
		h.handleDetectHTTP(w, r)
		return
	}
	_url := h.ServingAddr
	proxy := httputil.NewSingleHostReverseProxy(_url)

	// logger.Printf("Out:\t%v %s %s\n", _url, r.Method, r.URL)
	proxy.ServeHTTP(w, r)
}
