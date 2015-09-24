package manager

import (
	"github.com/samalba/dockerclient"
	"log"
	"os"
	"time"
	"crypto/rand"
	"fmt"
	"github.com/odewahn/swarm-manager/models"
	"github.com/odewahn/swarm-manager/db"
)

var (
	docker *dockerclient.DockerClient
	initialized bool
)

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
func Start(c *models.Container) {
	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}
	c.Hostname = getHostName()
  log.Println("Starting container based on ", c.Image)
	db.SaveContainer(c)	//Save startup state in the db
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
	db.SaveContainer(c)	//Save state in the db


}

// Kill a container
func Kill(c *models.Container) {
	if !initialized {
			log.Fatal("Package not initialized.  Call .Init() function.")
	}
	c.Status = "DELETING"
	db.SaveContainer(c)	//Save startup state in the db
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
	c.Status = "REMOVED"
	db.SaveContainer(c)	//Save startup state in the db

}


func NoOp(c *models.Container) {
	c.Hostname = getHostName()
	fmt.Println(c.Serialize())
	db.SaveContainer(c)
	time.Sleep(20 * time.Second)
  c.Status = "READY"
	db.SaveContainer(c)
	fmt.Println(c.Serialize())
}
