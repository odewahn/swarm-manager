package manager

import (
	"github.com/samalba/dockerclient"
	"log"
	"os"
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
	initialized bool
)

// Initialize the docker connection
func Init() {

	tlsConfig, err :=  getTLSConfig(os.Getenv("SWARM_CREDS_DIR"))
	if err != nil {
		log.Fatal("Could not create TLS certificate.")
	}

	docker, err = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)

	if err != nil {
		log.Fatal("Error initializing docker: ", err)
	}
	log.Println("Swarm connection inialized")

	initialized = true
}

// Start a container
func (c *Container) Start() {

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

  log.Println("Starting container based on ", c.Image)

  // Create the container
	containerConfig := &dockerclient.ContainerConfig{
		Image: c.Image,
		Cmd:   []string{"/bin/sh", "-c", "ipython notebook --ip=0.0.0.0 --no-browser"},
		ExposedPorts: map[string]struct{}{
			"8888/tcp": {},
		},
		Hostname: generateHostName("ANIMAL"),
		Domainname: c.Domainname,
	}

	containerId, err := docker.CreateContainer(containerConfig, containerConfig.Hostname)
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

	log.Println("Started container ", containerConfig.Hostname)

}

// Kill a container
func (c *Container) Kill() {

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

  log.Println("Stopping container ", c.ContainerId)
  err := docker.StopContainer(c.ContainerId, 5)
  if err != nil {
    log.Println("Could not kill container", c.ContainerId)
  }
	log.Println("Removing container ", c.ContainerId)
  docker.RemoveContainer(c.ContainerId, true, true)
  if err != nil {
    log.Println("Could not remove container ", c.ContainerId)
  }
	log.Println("Removed container ", c.ContainerId)

}
