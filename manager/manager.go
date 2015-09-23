package manager

import (
	"github.com/samalba/dockerclient"
	"log"
	"os"
	"time"
)

type Container struct {
  Hostname string `json:"hostname"`
  Domainname string `json:"domainname"`
  Image string `json:"image"`
  Url string `json:"url"`
  ContainerId string
	Status string
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
	// Setup the docker host
	docker, err = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)
	if err != nil {
		log.Fatal("Error initializing docker: ", err)
	}
	log.Println("Swarm connection inialized")
	initialized = true
}

// Start a container
func (c *Container) Start(status chan string) {

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

	status <- "STARTING"

	c.Hostname = generateHostName("ANIMAL")

  log.Println("Starting container based on ", c.Image)
  // Create the container
	containerConfig := &dockerclient.ContainerConfig{
		Image: c.Image,
		Cmd:   []string{"/bin/sh", "-c", "ipython notebook --ip=0.0.0.0 --no-browser"},
		ExposedPorts: map[string]struct{}{
			"8888/tcp": {},
		},
		Hostname: c.Hostname,
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

	status <- "DONE"

}

// Kill a container
func (c *Container) Kill(status chan string) {

	status <- "STARTING"

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

  log.Println("Stopping container ", c.Hostname)

  err := docker.StopContainer(c.ContainerId, 5)
  if err != nil {
    log.Println("Could not kill container", c.Hostname)
  }

	log.Println("Removing container ", c.Hostname)
  docker.RemoveContainer(c.Hostname, true, true)
  if err != nil {
    log.Println("Could not remove container ", c.Hostname)
  }

	status <- "DONE"

	log.Println("Removed container ", c.Hostname)

}


func (c *Container) NoOp(status chan string) {
	status <- "WORKING"
	time.Sleep(3 * time.Second)
	status <- "DONE"
}
