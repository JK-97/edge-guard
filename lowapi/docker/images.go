package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jxcore/lowapi/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

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
		logger.Error(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := c.cli.ContainerStop(c.ctx, container.ID, nil); err != nil {
			logger.Error(err)
		}
		fmt.Println("Success")
	}
}

//ImagesList is
func (c *DockerObj) ImagesList() ([]*DockerImageResources, error) {
	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		logger.Error(err)
	}
	s := make([]*DockerImageResources, 0)
	for _, image := range images {
		image_info := &DockerImageResources{Tag: image.RepoTags, ImageID: image.ID, Lable: image.Labels}
		s = append(s, image_info)

	}

	return s, err
}

//ImageAllRemove remove all images
func (c *DockerObj) ImageAllRemove() {
	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		logger.Error(err)
	}
	for _, image := range images {
		_, err := c.cli.ImageRemove(c.ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			logger.Error(err)
		}
		logger.Info("has delete image : " + image.ID)
	}

}

// DockerRestore 恢复 docker 镜像
func (c *DockerObj) DockerRestore() error {
	// get docker_desc from desc.json
	logger.Info("Reading from ", DockerDesc)
	data, err := ioutil.ReadFile(DockerDesc)
	if err != nil {
		return err
	}
	m := make(map[string]map[string]string)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	for filename, info := range m {
		logger.Info("Restore Image: ", info["repo"])
		path := filepath.Join(DockerRestorePath, filename+".tar")
		fileObj, err := os.Open(path)
		if err != nil {
			path = strings.ReplaceAll(path, ":", "_")
		}
		fileObj, err = os.Open(path)
		if err == nil {
			defer fileObj.Close()
			_, err := c.cli.ImageLoad(c.ctx, fileObj, true)
			if err != nil {
				logger.Error(err)
			}
			c.cli.ImageTag(c.ctx, info["id"], info["repo"])
		} else {
			logger.Error("Import Failed", filename, err)
		}

	}
	return nil
}

func (c *DockerObj) Removebyimage(imagename string) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {

		logger.Error("docker", err)
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
