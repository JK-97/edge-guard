package docker

import (
    "github.com/docker/docker/client"
    "sync"
    "context"
)

//DockerObj is
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
