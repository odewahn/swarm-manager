package container

import (
	"github.com/samalba/dockerclient"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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


// Reads the specified file into a byte array
func fetchFile(path, fn string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/%s", path, fn)
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error loading %s : %s", fn, err)
	}
	return out, err
}

// Loads and returns a tls config settings for the swarm.  Lightly adapted version of:
//     https://github.com/ehazlett/interlock/blob/master/interlock/main.go#L14-L32
func GetTLSConfig(certsDir string) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config

	caCert, err := fetchFile(certsDir, os.Getenv("SWARM_CA"))
	cert, err := fetchFile(certsDir, os.Getenv("SWARM_CERT"))
	key, err := fetchFile(certsDir, os.Getenv("SWARM_KEY"))

	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}

	return &tlsConfig, nil
}


// Initialize the docker connection
func Init(d *dockerclient.DockerClient) {
	docker = d
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
