package docker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"jxcore/lowapi/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

//BuildImage build image
func BuildImage(fileReader io.Reader, imageName string) error {

	buildOptions := types.ImageBuildOptions{
		// Dockerfile: "Dockerfile", // optional, is the default
		Tags: []string{imageName},
		// Labels: map[string]string{
		// 	project: "project",
		// },
	}

	output, err := dockerObj.cli.ImageBuild(context.Background(), fileReader, buildOptions)
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

//ImagesList is
func ImagesList() ([]*DockerImageResources, error) {
	images, err := dockerObj.cli.ImageList(dockerObj.ctx, types.ImageListOptions{})
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
func ImageAllRemove() {
	images, err := dockerObj.cli.ImageList(dockerObj.ctx, types.ImageListOptions{})
	if err != nil {
		logger.Error(err)
	}
	for _, image := range images {
		_, err := dockerObj.cli.ImageRemove(dockerObj.ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			logger.Error(err)
		}
		logger.Info("has delete image : " + image.ID)
	}

}

func LoadImage(fileReader io.Reader) error {
	_, err := dockerObj.cli.ImageLoad(dockerObj.ctx, fileReader, true)
	return err
	// dockerObj.cli.ImageTag(dockerObj.ctx, info["id"], info["repo"])
}

// DockerRestore 恢复 docker 镜像
func DockerRestore() error {
	// get docker_desc from desdockerObj.json
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
			_, err := dockerObj.cli.ImageLoad(dockerObj.ctx, fileObj, true)
			if err != nil {
				logger.Error(err)
			}
			dockerObj.cli.ImageTag(dockerObj.ctx, info["id"], info["repo"])
		} else {
			logger.Error("Import Failed", filename, err)
		}

	}
	return nil
}

func Removebyimage(imagename string) {
	containers, err := dockerObj.cli.ContainerList(dockerObj.ctx, types.ContainerListOptions{})
	if err != nil {

		logger.Error("docker", err)
	}

	for _, perone := range containers {
		if perone.Image == imagename {
			thecontanierid := perone.ID
			//theimageid := perone.ImageID
			dockerObj.cli.ContainerRemove(dockerObj.ctx, thecontanierid, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
			//dockerObj.cli.ImageRemove(dockerObj.ctx, theimageid, types.ImageRemoveOptions{Force: true, PruneChildren: true})
			break
		}
	}
}
