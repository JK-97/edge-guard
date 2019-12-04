package docker

import (
	"encoding/json"
	"fmt"
	log "jxcore/lowapi/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

//var ctx = context.Background()

//NewClient return a docker client
func NewClient() (dockerobj DockerObj) {
	var err error
	dockerobj.Messages = make(chan string)
	dockerobj.Mu = new(sync.Mutex)
	dockerobj.dockersingle = nil
	dockerobj.ctx = context.Background()
	dockerobj.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
	}
	return dockerobj
}

//var dockersingle *client.Client
//var mu sync.Mutex

func (c *DockerObj) GetClient() *client.Client {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if c.dockersingle == nil {
		c.dockersingle, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	} // unnecessary locking if instance already created

	return c.dockersingle
}

//BuildImage build image
func (c *DockerObj) BuildImage(tarFile, project, imageName string) error {
	dockerBuildContext, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{imageName},
		Labels: map[string]string{
			project: "project",
		},
	}

	output, err := c.cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "error") {
		return fmt.Errorf("build image to docker error")
	}

	return nil
}

//StopContainer return docker images list
func (c *DockerObj) StopContainer() {

	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := c.cli.ContainerStop(c.ctx, container.ID, nil); err != nil {
			log.Error(err)
		}
		fmt.Println("Success")
	}
}

//ImagesList is
func (c *DockerObj) ImagesList() ([]*DockerImageResources, error) {
	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		log.Error(err)
	}
	//for per_image := range images{
	//	per_image_info:= &schema.DockerInternalResources{Tag:per_image.tag}
	//}
	s := make([]*DockerImageResources, 0)
	for _, image := range images {
		//str, _ := json.Marshal(image)
		//fmt.Printf("%s\n", str)
		image_info := &DockerImageResources{Tag: image.RepoTags, ImageID: image.ID, Lable: image.Labels}
		s = append(s, image_info)

	}

	//dockerresp := &DockerImageResp{Status: "runing", ProcessPercentage: 0.2, DockerInternalResources: s}
	return s, err
}

func (c *DockerObj) ContainerList() ([]*DockerContainerResources, error) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}

	s := make([]*DockerContainerResources, 0)

	for _, container := range containers {
		contaner_info := &DockerContainerResources{Status: container.Status, ContainerID: container.ID, Image: container.Image, Lable: container.Labels}
		s = append(s, contaner_info)
	}
	//dockerresp := &DockerContainerResp{Status: "runing", ProcessPercentage: 0.2, DockerInternalResources: s}
	return s, err
}

//ContainerAllRemove :it will remove all container
func (c *DockerObj) ContainerAllRemove() {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, container := range containers {
		err := c.cli.ContainerRemove(c.ctx, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		if err != nil {
			log.Error(err)
		}
		log.Info("has delete coontainer : " + container.ID)
	}

}

//ImageAllRemove remove all images
func (c *DockerObj) ImageAllRemove() {
	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		log.Error(err)
	}
	// log.Info(images)
	for _, image := range images {
		_, err := c.cli.ImageRemove(c.ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			log.Error(err)
		}
		log.Info("has delete image : " + image.ID)
	}

}

// DockerRestore 恢复 docker 镜像
func (c *DockerObj) DockerRestore() error {
	// get docker_desc from desc.json
	log.Info("Reading from ", DockerDesc)
	data, err := ioutil.ReadFile(DockerDesc)
	if err != nil {
		// log.Error(err)
		return err
	}
	m := make(map[string]map[string]string)
	err = json.Unmarshal(data, &m)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	for filename, info := range m {
		// model.BufioRead("")
		log.Info("Restore Image: ", info["repo"])
		path := filepath.Join(DockerRestorePath, filename+".tar")
		fileObj, err := os.Open(path)
		if err != nil {
			path = strings.ReplaceAll(path, ":", "_")
		}
		fileObj, err = os.Open(path)
		if err == nil {
			defer fileObj.Close()
			_, err := c.cli.ImageLoad(c.ctx, fileObj, true)
			// _, err := c.cli.ImageImport(c.ctx, types.ImageImportSource{Source: fileObj, SourceName: "-"}, info["repo"], types.ImageImportOptions{Tag: info["tag"]})
			if err != nil {
				log.Error(err)
			}
			c.cli.ImageTag(c.ctx, info["id"], info["repo"])
		} else {
			log.Error("Import Failed", filename, err)
		}

	}
	return nil
}

func (c *DockerObj) Removebyimage(imagename string) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {

		log.Error("docker", err)
	}

	for _, perone := range containers {
		if perone.Image == imagename {
			thecontanierid := perone.ID
			//theimageid := perone.ImageID
			c.cli.ContainerRemove(c.ctx, thecontanierid, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
			//c.cli.ImageRemove(c.ctx, theimageid, types.ImageRemoveOptions{Force: true, PruneChildren: true})
			break
		}
	}
}
