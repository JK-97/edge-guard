package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	log "jxcore/lowapi/logger"

	yaml "gopkg.in/yaml.v2"
)

// DockerServeMode docker 代理模式
type DockerServeMode uint

const composeDirPrefix = "compose-"

const (
	// DockerServeModeUnix unix
	DockerServeModeUnix DockerServeMode = iota
	// DockerServeModeTCP tcp
	DockerServeModeTCP
)

// DockerProxy Dockerd 代理
type DockerProxy struct {
	Domain    string
	ServeMode DockerServeMode
	handler   http.Handler
}

// DockerComposeAgent docker-compose 代理
type DockerComposeAgent struct {
	ComposeBinary string
	BaseDir       string
}

// DockerComposeRequest docker compose 相关请求
type DockerComposeRequest struct {
	File string                 `json:"file"`
	Dir  string                 `json:"dir"`
	Path string                 `json:"path"`
	Args []string               `json:"args"`
	Yaml map[string]interface{} `json:"yaml"`
}

// GetAllPaths 获取路径
func (req *DockerComposeRequest) GetAllPaths() (dir, file, path string) {
	dir, file, path = req.Dir, req.File, req.Path

	if dir == "" {
		dir = filepath.Dir(path)
	}
	if file == "" {
		file = filepath.Base(path)
	}
	if path == "" {
		path = filepath.Join(dir, file)
	}

	return
}

type unixProxyHandler struct {
	// unixDial func(network, addr string) (net.Conn, error)
	client *http.Client
}

func (h *unixProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := h.client

	req, _ := http.NewRequest(
		r.Method,
		fmt.Sprintf("http://127.0.0.1%s?%s", r.URL.Path, r.URL.RawQuery),
		r.Body,
	)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	resp.Write(w)
}

// NewDockerProxy 获取新的 DockerProxy 实例
func NewDockerProxy(domain string) *DockerProxy {
	if domain == "" {
		domain = "http://127.0.0.1:2375"
	}

	uri, err := url.Parse(domain)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	if uri.Scheme == "unix" {
		var d net.Dialer
		handler := unixProxyHandler{
			client: &http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return d.DialContext(ctx, "unix", uri.Path)
					},
				},
			},
		}
		return &DockerProxy{
			Domain:    domain,
			ServeMode: DockerServeModeUnix,
			handler:   &handler,
		}
	}

	return &DockerProxy{
		Domain:    domain,
		ServeMode: DockerServeModeTCP,
		handler:   httputil.NewSingleHostReverseProxy(uri),
	}
}

// NewDockerComposeAgent 获取 DockerComposeAgent 实例
// bin docker-compose 文件路径
// base docker-compose 基础路径
func NewDockerComposeAgent(bin, base string) *DockerComposeAgent {
	if bin == "" {
		bin = "docker-compose"
	}
	if base == "" {
		base = "/data/compose"
	}
	err := os.MkdirAll(base, 0755)
	if err != nil {
		log.Error(err)
	}
	return &DockerComposeAgent{bin, base}
}

func (h *DockerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.handler == nil {
		panic("not impl")
	}

	h.handler.ServeHTTP(w, r)
}

// filterInvalidComposeSection 过滤掉 docker compose 文件不支持的 key
func filterInvalidComposeSection(data *map[string]interface{}) {
	// Valid top-level sections for this Compose file are: version, services, networks, volumes, and extensions starting with "x-".
	keys := make([]string, 0, len(*data))

	for k := range *data {
		switch k {
		case "version", "services", "networks", "volumes":
		default:
			if !strings.HasPrefix(k, "x-") {
				keys = append(keys, k)
			}
		}
	}
	for _, key := range keys {
		delete(*data, key)
	}

}

func jsonToYaml(buffer []byte) (result []byte, err error) {
	data := make(map[string]interface{})

	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return nil, err
	}

	filterInvalidComposeSection(&data)

	result, err = yaml.Marshal(data)

	return
}

func (h *DockerComposeAgent) readJSONBody(w http.ResponseWriter, r *http.Request) (buffer []byte, err error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		ErrorWithCode(w, http.StatusBadRequest)
		return
	}
	buffer, err = ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorWithCode(w, http.StatusBadRequest)
		return
	}
	return
}

func captureCommandResult(cmd *exec.Cmd, w http.ResponseWriter) (output []byte, err error) {

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		code := cmd.ProcessState.ExitCode()
		Error(w, fmt.Sprintf("Code: %d, GoErr: %s, Compose: %s", code, err.Error(), string(output)), http.StatusInternalServerError)
	}
	return
}

// cleanUp 清理 Docker Compose 目录
func (h *DockerComposeAgent) cleanUp(r *http.Request, dir, file string) {
	<-r.Context().Done()
	log.Infof("EndRequest :\t%s %s %s", r.RemoteAddr, r.Method, r.URL)
	var requestFailed bool
	if r.Response != nil {
		resp := r.Response
		log.Infof("Response Code:\t%d %s %s", resp.StatusCode, resp.Status)
		if resp.StatusCode >= 400 {
			requestFailed = true
		}
	} else {
		log.Info("No Response")
		requestFailed = true
	}

	if r.Context().Err() != nil {
		requestFailed = true
		log.Infof("Error: %s", r.Context().Err())
	}

	if requestFailed {
		if _, err := os.Stat(dir); err != nil {
			cmd := exec.Command(h.ComposeBinary, "-f", file, "down")
			cmd.Dir = dir

			cmd.Run()
			os.RemoveAll(dir)
		}
	}
}

// dockerComposeUp 创建 services
func (h *DockerComposeAgent) dockerComposeUp(w http.ResponseWriter, r *http.Request) {

	buffer, err := h.readJSONBody(w, r)
	if err != nil {
		return
	}

	yml, err := jsonToYaml(buffer)

	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dir string
	file := "docker-compose.yml"
	composeID := r.Header.Get("JX-Compose-Id")
	if composeID == "" {
		// 旧版 Deploy Engine 请求
		dir, err = ioutil.TempDir("/data/compose", composeDirPrefix)
		if err != nil {
			log.Println(err)
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		os.Chmod(dir, 0755)
	} else {
		dir = filepath.Join("/data/compose", composeID)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		} else {
			// 卸载旧应用
			cmd := exec.Command(h.ComposeBinary, "-f", file, "down")
			cmd.Dir = dir
			cmd.Run()
		}
	}

	err = ioutil.WriteFile(filepath.Join(dir, file), yml, 0666)
	if err != nil {
		log.Println(err)
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 请求中断时，删除 目录和 Docker 容器
	go h.cleanUp(r, dir, file)

	cmd := exec.Command(h.ComposeBinary, "-f", file, "up", "-d", "--remove-orphans")
	cmd.Dir = dir

	if _, err = captureCommandResult(cmd, w); err != nil {
		// os.RemoveAll(dir)
		return
	}

	data := make(map[string]interface{})
	data["file"] = file
	data["dir"] = dir
	data["path"] = filepath.Join(dir, file)
	WriteResult(w, NewAPIResult(&data))

	log.Printf("Create Services: %s %s\n", dir, file)
}

// 更新 yaml 文件
func updateYaml(w http.ResponseWriter, req *DockerComposeRequest, path string) error {
	if req.Yaml != nil && len(req.Yaml) > 0 {
		// 更新 Yaml 文件
		filterInvalidComposeSection(&(req.Yaml))
		yml, err := yaml.Marshal(req.Yaml)

		if err != nil {
			Error(w, err.Error(), http.StatusBadRequest)
			return err
		}

		err = ioutil.WriteFile(path, yml, 0666)
		if err != nil {
			log.Println(err)
			Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

// dockerComposeUpdate 更新 services
func (h *DockerComposeAgent) dockerComposeUpdate(w http.ResponseWriter, r *http.Request) {
	buffer, err := h.readJSONBody(w, r)
	if err != nil {
		return
	}

	req := DockerComposeRequest{}

	if err := json.Unmarshal(buffer, &req); err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, file, path := req.GetAllPaths()

	if err = updateYaml(w, &req, path); err != nil {
		return
	}

	cmd := exec.Command(h.ComposeBinary, "-f", file, "up", "-d")
	cmd.Dir = dir

	if _, err = captureCommandResult(cmd, w); err != nil {
		return
	}

	data := make(map[string]interface{})
	data["file"] = file
	data["dir"] = dir
	WriteResult(w, NewAPIResult(&data))

	log.Printf("Update Services: %s %s\n", dir, file)
}

func (h *DockerComposeAgent) handleComposeCommand(w http.ResponseWriter, r *http.Request, command string, withOutput bool) {
	buffer, err := h.readJSONBody(w, r)
	if err != nil {
		return
	}

	req := DockerComposeRequest{}

	if err := json.Unmarshal(buffer, &req); err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, file, path := req.GetAllPaths()

	cmd := exec.Command(h.ComposeBinary, "-f", file, command)

	if req.Args != nil {
		cmd.Args = append(cmd.Args, req.Args...)
	}

	cmd.Dir = dir

	data := make(map[string]interface{})
	data["file"] = file
	data["dir"] = dir
	data["path"] = path

	output, err := captureCommandResult(cmd, w)
	if err != nil {
		return
	}
	if withOutput {
		data["result"] = string(output)
	}

	WriteResult(w, NewAPIResult(&data))
}

// DockerComposeCommand 执行 docker compose 命令，并忽略输出
func (h *DockerComposeAgent) DockerComposeCommand(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	cmd := data["command"]

	h.handleComposeCommand(w, r, cmd, false)
}

// DockerComposeWithOutput 执行 docker compose 命令，并捕获输出
func (h *DockerComposeAgent) DockerComposeWithOutput(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	cmd := data["command"]

	h.handleComposeCommand(w, r, cmd, true)
}

// dockerComposeDown 删除 services
func (h *DockerComposeAgent) dockerComposeDown(w http.ResponseWriter, r *http.Request) {
	req := DockerComposeRequest{}
	contentType := r.Header.Get("Content-Type")
	var err error
	if strings.Contains(contentType, "json") {
		buffer, err := h.readJSONBody(w, r)
		if err != nil {
			return
		}
		if err := json.Unmarshal(buffer, &req); err != nil {
			Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		query := r.URL.Query()
		req.File = query.Get("file")
		req.Dir = query.Get("dir")
		req.Path = query.Get("path")
	}

	dir, file, path := req.GetAllPaths()

	cmd := exec.Command(h.ComposeBinary, "-f", file, "down", "--remove-orphans")
	cmd.Dir = dir

	if _, err = captureCommandResult(cmd, w); err != nil {
		return
	}

	data := make(map[string]interface{})
	data["file"] = file
	data["dir"] = dir
	data["path"] = path
	WriteResult(w, NewAPIResult(&data))

	os.RemoveAll(dir)
	log.Printf("Remove Services: %s %s\n", dir, file)
}

func checkResult(c chan<- *exec.Cmd, cmd *exec.Cmd) {
	err := cmd.Run()
	if err == nil {
		c <- nil
		os.RemoveAll(cmd.Dir)
	} else {
		c <- cmd
	}
}

// DockerComposeTruncate 清除通过 Gateway 创建的所有 services
func (h *DockerComposeAgent) DockerComposeTruncate(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(h.BaseDir)
	if err != nil {
		WriteSucess(w)
		return
	}
	dirCount := len(files)

	data := make(map[string]interface{}, dirCount)

	resultChan := make(chan *exec.Cmd, dirCount)

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), composeDirPrefix) {

			cmd := exec.Command(h.ComposeBinary, "-f", "docker-compose.yml", "down")
			cmd.Dir = filepath.Join(h.BaseDir, file.Name())

			b := new(bytes.Buffer)
			cmd.Stdout = b
			cmd.Stderr = b

			go checkResult(resultChan, cmd)
			data[cmd.Dir] = ""
		} else {
			dirCount--
		}
	}
	data["total"] = dirCount

	for index := 0; index < dirCount; index++ {
		cmd := <-resultChan
		if cmd != nil {
			b, _ := cmd.Stdout.(*bytes.Buffer)
			if b != nil {
				data[cmd.Dir] = fmt.Sprintf("Code: %d, Compose: %s", cmd.ProcessState.ExitCode(), b.String())
			} else {
				data[cmd.Dir] = fmt.Sprintf("Code: %d", cmd.ProcessState.ExitCode())
			}
		}
	}

	WriteData(w, &data)
}

func (h *DockerComposeAgent) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.dockerComposeUp(w, r)
		return
	case http.MethodPut:
		h.dockerComposeUpdate(w, r)
		return
	case http.MethodDelete:
		h.dockerComposeDown(w, r)
		return
	default:
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
}
