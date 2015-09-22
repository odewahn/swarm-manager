package container

import (
	"github.com/samalba/dockerclient"
  "log"
)

type Container struct {
  Hostname string `json:"hostname"`
  Domainname string `json:"domainname"`
  Image string `json:"image"`
  Url string `json:"url"`
  ContainerId string
}

var (
	docker *dockerclient.DockerClient
)

func Init(d *dockerclient.DockerClient) {
	docker = d
}

func (c *Container) Start() {
  log.Println("Starting container based on ", c.Image)

  // Create the container
	containerConfig := &dockerclient.ContainerConfig{
		Image: c.Image,
		Cmd:   []string{"/bin/sh", "-c", "ipython notebook --ip=0.0.0.0 --no-browser"},
		ExposedPorts: map[string]struct{}{
			"8888/tcp": {},
		},
		Hostname:   c.Hostname,
		Domainname: c.Domainname,
	}

	containerId, err := docker.CreateContainer(containerConfig, c.Hostname)
	if err != nil {
		log.Println(err)
	}

  c.ContainerId = containerId

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}

	err = docker.StartContainer(containerId, hostConfig)
	if err != nil {
		log.Println(err)
	}

	log.Println("Started container ", c.ContainerId)

}

func (c *Container) Kill() {

  log.Printf("Killing container ", c.ContainerId)
  err := docker.StopContainer(c.ContainerId, 5)
  if err != nil {
    log.Printf("Could not kill container %s", c.ContainerId)
  }
  docker.RemoveContainer(c.ContainerId,true,true)
  if err != nil {
    log.Printf("Could not remove container ", c.ContainerId)
  }

}
