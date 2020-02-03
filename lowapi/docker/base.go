package docker

import (
	"jxcore/lowapi/logger"
	"sync"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var dockerObj DockerObj

const (
	DockerRestorePath = "/restore/dockerimage/"
	DockerDesc        = "/restore/dockerimage/desc.json"

	daemonConfigPath = "/etc/docker/daemon.json"
)

type DockerObj struct {
	Messages     chan string
	Mu           *sync.Mutex
	dockersingle *client.Client
	ctx          context.Context
	cli          *client.Client
}

type DockerImageResp struct {
	Status                  string                  `json:"status"`
	DockerInternalResources []*DockerImageResources `json:"docker_internal_resources"`
	ProcessPercentage       float32                 `json:"process_percentage"`
}

//DockerContainerResp is docker container info
type DockerContainerResp struct {
	Status                  string                      `json:"status"`
	DockerInternalResources []*DockerContainerResources `json:"docker_internal_resources"`
	ProcessPercentage       float32                     `json:"process_percentage"`
}

//DockerImageResources is
type DockerImageResources struct {
	Tag     []string          `json:"tag"`
	ImageID string            `json:"image_id"`
	Lable   map[string]string `json:"label"`
}

//DockerImageResources is
type DockerContainerResources struct {
	Image       string            `json:"image"`
	ContainerID string            `json:"contaner_id"`
	Lable       map[string]string `json:"label"`
	Status      string            `json:"status"`
}

type DockerRestoreDesc struct {
	FileName string `json:"image"`
	Tag      string `json:"tag"`
	Repo     string `json:"repo"`
}

//NewClient return a docker client
func NewClient() (d DockerObj) {
	var err error
	d.Messages = make(chan string)
	d.Mu = new(sync.Mutex)
	d.dockersingle = nil
	d.ctx = context.Background()
	d.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error(err)
	}
	return d
}

//var mu sync.Mutex

func (c *DockerObj) GetClient() *client.Client {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if c.dockersingle == nil {
		c.dockersingle, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	} // unnecessary locking if instance already created

	return c.dockersingle
}

func init() {
	dockerObj = NewClient()
}
