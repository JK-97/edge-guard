package gateway

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"jxcore/lowapi/logger"

	"jxcore/gateway/option"
	"jxcore/gateway/serve"
	"jxcore/gateway/store"
	"jxcore/web"

	"github.com/gorilla/mux"
)

// ServerOptions 服务配置
var ServerOptions option.ServerConfig
var defaultStore store.Store

func makeRouter() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = serve.NewNotFoundHandler()
	r.Use(simpleMw)
	r.Use(recoverMiddleware)

	config := ServerOptions.ProxyServerConfig
	proxy := serve.CreateReverseProxyFromOpion(&config)

	proxyRouter := r.PathPrefix("/internalapi/").Name("reverse").Subrouter()
	proxyRouter.PathPrefix("/").Handler(proxy)

	appendMessageQueueRouter(r)
	appendAiServingRouter(r)
	appendContainerRouter(r)
	appendPowerRouter(r)
	appendConfigAgentRouter(r)
	appendDeviceRouter(r)
	appendDatabaseRouter(r)

	appendTSDBRouter(r)

	appendFileRouter(r)

	r.HandleFunc("/api/v1/nodeinfo", serve.HandleGetNodeInfo)

	if ServerOptions.EnableDynamicService {
		appendDynamicServiceRouter(r)
	}

	return r
}

func simpleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("In :\t%s %s %s", r.RemoteAddr, r.Method, r.URL)
		start := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		cost := time.Since(start)
		resp := r.Response
		if resp != nil {
			logger.Infof("End:\t%s %s %s %dms %s", r.RemoteAddr, r.Method, r.URL, cost.Milliseconds(), resp.Status)
		} else {
			logger.Infof("End:\t%s %s %s %dms", r.RemoteAddr, r.Method, r.URL, cost.Milliseconds())
		}
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("Unknown panic")
				}
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// appendMessageQueueRouter 添加 MQ 服务需要的路由
func appendMessageQueueRouter(r *mux.Router) {
	prefix := "/api/v1/mq"
	msgQueueRouter := r.PathPrefix(prefix).Subrouter()

	mqHandler := serve.NewMesageQueueInternalHandler(ServerOptions.MessageQueueConfig)
	if defaultStore != nil {
		mqHandler.SetStore(defaultStore)
	}

	stripedMQHandler := http.StripPrefix(prefix, mqHandler)
	msgQueueRouter.Path("/grant/{topic}").Handler(stripedMQHandler)
	msgQueueRouter.PathPrefix("/").Handler(stripedMQHandler)
}

// appendAiServingRouter 添加 AI Serving 需要的路由
func appendAiServingRouter(r *mux.Router) {
	prefix := "/api/v1/ai"
	aiRouter := r.PathPrefix(prefix).Subrouter()
	aiHandler := serve.NewAiServingHandler(ServerOptions.AiServingURL)
	stripedHandler := http.StripPrefix(prefix, aiHandler)
	aiRouter.PathPrefix("/").Handler(stripedHandler)
}

// appendContainerRouter 容器相关的路由
func appendContainerRouter(r *mux.Router) {
	prefix := "/api/v1/docker"
	router := r.PathPrefix(prefix).Subrouter()
	handlerImpl := serve.NewDockerProxy(ServerOptions.DockerDomin)
	handler := http.StripPrefix(prefix, handlerImpl)
	router.PathPrefix("/").Handler(handler)

	prefix = "/api/v1/compose"
	base := "/data/compose"
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}

	router = r.PathPrefix(prefix).Subrouter()
	handlerImpl2 := serve.NewDockerComposeAgent(ServerOptions.ComposeBinary, ServerOptions.ComposeBaseDir)

	router.HandleFunc(
		"/{command:create|down|kill|pause|port|pull|push|restart|rm|scale|start|stop|unpause|up}",
		handlerImpl2.DockerComposeCommand)
	router.HandleFunc(
		"/{command:config|events|exec|images|logs|ps|run|scale|top|version}",
		handlerImpl2.DockerComposeWithOutput)
	router.HandleFunc("/truncate", handlerImpl2.DockerComposeTruncate)

	handler = http.StripPrefix(prefix, handlerImpl2)
	router.PathPrefix("").Handler(handler)
}

// appendDynamicServiceRouter 添加动态服务相关的路由
func appendDynamicServiceRouter(r *mux.Router) {

	dynamicProxy := serve.NewDynamicServiceContainer()
	r.PathPrefix("/api/dynamic").Name("service-registration").Handler(dynamicProxy)
	r.PathPrefix("/dynamicapi/").Name("dynamic").Handler(http.StripPrefix("/dynamicapi", dynamicProxy.DynamicHandler()))
	// s := store.NewLevelDBStore(filepath.Join(*workingDir, "gateway.db"))
	s := defaultStore
	if defaultStore == nil {
		panic(errors.New("No Store Configed"))
	}
	dynamicProxy.Recovery(s)

	(*dynamicProxy).OnChange = func() {
		dynamicProxy.Store(s)
	}
}

// appendPowerRouter 添加电源相关的路由
func appendPowerRouter(r *mux.Router) {
	r.Path("/api/v1/power/off").Methods(http.MethodPost).HandlerFunc(serve.PowerOffHTTP)
	r.Path("/api/v1/power/mode").HandlerFunc(serve.StartUpMode)
}

// appendConfigAgentRouter 添加 Config Agent 相关的路由
func appendConfigAgentRouter(r *mux.Router) {
	prefix := "/api/v1/config"

	c := serve.NewConfigAgentHandler(ServerOptions.ConfigAgent)
	router := r.PathPrefix(prefix).Subrouter()

	router.PathPrefix("/ws").HandlerFunc(c.ServeWebsocket)
	router.PathPrefix("/").Handler(http.StripPrefix(prefix, c))
}

// appendDeviceRouter 添加 Device 相关的路由
func appendDeviceRouter(r *mux.Router) {
	// url: /api/v1/device/command(或者其他edgex模块名)/转发的url
	prefix := "/api/v1/device/"

	c := serve.NewDeviceHandler(ServerOptions.Device)
	router := r.PathPrefix(prefix).Subrouter()

	router.PathPrefix("/").Handler(http.StripPrefix(prefix, c))
}

// appendDatabaseRouter 添加 Device 相关的路由
func appendDatabaseRouter(r *mux.Router) {
	prefix := "/api/v1/db"

	c := serve.NewDatabaseHandler(ServerOptions.Database)
	router := r.PathPrefix(prefix).Subrouter()

	router.PathPrefix("/").Handler(http.StripPrefix(prefix, c))
}

// appendTSDBRouter 添加时序数据库 相关的路由
func appendTSDBRouter(r *mux.Router) {
	prefix := "/api/v1/tsdb"

	c := serve.NewTSDBHandler(ServerOptions.TimeSeries)
	c.SetStore(defaultStore)
	router := r.PathPrefix(prefix).Subrouter()

	router.PathPrefix("/").Handler(http.StripPrefix(prefix, c))
}

// appendFileRouter 添加文件操作相关的路由
func appendFileRouter(r *mux.Router) {
	prefix := "/api/v1/file"

	c := &serve.FileHandler{}
	router := r.PathPrefix(prefix).Subrouter()

	router.PathPrefix("/").Handler(http.StripPrefix(prefix, c))
}

// listenUnix 监听 Unix Socket
func listenUnix(router http.Handler) {
	addr := net.UnixAddr{Name: ServerOptions.SocketAddr, Net: "unix"}

	if err := os.Remove(addr.Name); err != nil {
		if !os.IsNotExist(err) {
			logger.Info(err)
		}
	}

	ln, err := net.ListenUnix("unix", &addr)

	if err != nil {
		logger.Info(err)
		return
	}
	if ServerOptions.SocketMode > 0 {
		os.Chmod(addr.Name, os.FileMode(ServerOptions.SocketMode))
	}

	defer os.Remove(addr.Name)
	logger.Info("Listen on: unix: ", addr.Name)
	http.Serve(ln, router)
}

// ServeGateway 启动 Gateway
func ServeGateway(ctx context.Context, graceful time.Duration) error {
	Setup()
	router := makeRouter()

	// 监听 Unix Socket
	if ServerOptions.SocketAddr != "" {
		go listenUnix(router)
	}

	logger.Info("Listen on: http://", ServerOptions.Addr)
	return web.Serve(ctx, ServerOptions.Addr, router, graceful)
}
