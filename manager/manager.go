package manager

import (
	"github.com/samalba/dockerclient"
	"log"
	"os"
	"time"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

type Container struct {
  Hostname string `json:"hostname"`
  Domainname string `json:"domainname"`
  Image string `json:"image"`
  ContainerId string
	Status string
}

var (
	docker *dockerclient.DockerClient
	initialized bool
)

// Serializes a container as a string
func (c *Container) Serialize() (string) {
	out, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

// From https://www.socketloop.com/tutorials/golang-how-to-generate-random-string
func getHostName() string {
	dictionary := "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 12)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}


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
func (c *Container) Start() {

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

	c.Hostname = getHostName()

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

	c.Status = "ACTIVE"

}

// Kill a container
func (c *Container) Kill() {

	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}

  err := docker.StopContainer(c.ContainerId, 5)
  if err != nil {
    log.Println("Could not kill container", c.Hostname)
  }

	log.Println("Removing container ", c.Hostname)
  docker.RemoveContainer(c.Hostname, true, true)
  if err != nil {
    log.Println("Could not remove container ", c.Hostname)
  }

	log.Println("Removed container ", c.Hostname)

	c.Status = "DELETED"
}


func (c *Container) NoOp() {
	fmt.Println(c.Serialize())
	time.Sleep(3 * time.Second)
  c.Status = "READY"
	fmt.Println(c.Serialize())

}
