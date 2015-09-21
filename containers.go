package main

import (
	"github.com/samalba/dockerclient"
  "log"
)

type Containers struct {
  Client dockerclient.DockerClient
  Hostname string `json:"hostname"`
  Domainname string `json:"domainname"`
  Image string `json:"image"`
  Url string `json:"url"`
  ContainerId string
}


func (c *Containers) Start() (error) {
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

	containerId, err := c.Client.CreateContainer(containerConfig, c.Hostname)
	if err != nil {
		log.Println(err)
	}

  c.ContainerId = containerId

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}
	err = c.Client.StartContainer(containerId, hostConfig)
	if err != nil {
		log.Println(err)
	}
	return err
}


func (c *Containers) Kill() {

  log.Printf("Killing container ", c.ContainerId)
  err := c.Client.StopContainer(c.ContainerId, 5)
  if err != nil {
    log.Printf("Could not kill container %s", c.ContainerId)
  }
  c.Client.RemoveContainer(c.ContainerId,true,true)
  if err != nil {
    log.Printf("Could not remove container ", c.ContainerId)
  }

}
