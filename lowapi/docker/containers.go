package docker

import (
	"jxcore/lowapi/logger"

	"github.com/docker/docker/api/types"
)

func (c *DockerObj) ContainerList() ([]*DockerContainerResources, error) {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
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
func (c *DockerObj) ContainerAllRemove() {
	containers, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(err)
	}
	for _, container := range containers {
		err := c.cli.ContainerRemove(c.ctx, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		if err != nil {
			logger.Error(err)
		}
		logger.Info("has delete coontainer : " + container.ID)
	}

}
