package docker

import (
	"github.com/JK-97/edge-guard/lowapi/logger"

	"github.com/docker/docker/api/types"
)

func ContainerList() ([]*DockerContainerResources, error) {
	containers, err := dockerObj.cli.ContainerList(dockerObj.ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(err)
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
func ContainerAllRemove() {
	containers, err := dockerObj.cli.ContainerList(dockerObj.ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(err)
	}
	for _, container := range containers {
		err := dockerObj.cli.ContainerRemove(dockerObj.ctx, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		if err != nil {
			logger.Error(err)
		}
		logger.Info("has delete coontainer : " + container.ID)
	}

}

//StopContainer return docker images list
func StopContainer() {

	containers, err := dockerObj.cli.ContainerList(dockerObj.ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(err)
	}

	for _, container := range containers {
		logger.Info("Stopping container ", container.ID[:10], "... ")
		if err := dockerObj.cli.ContainerStop(dockerObj.ctx, container.ID, nil); err != nil {
			logger.Error(err)
		}
		logger.Info("Success")
	}
}
