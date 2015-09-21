package main

import (
  "fmt"
  "github.com/joho/godotenv"
  "github.com/samalba/dockerclient"
  "log"
  "os"
)

func main() {

  // Load the environment variables we need
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  tlsConfig, err := getTLSConfig(os.Getenv("SWARM_CREDS_DIR"))
  if err != nil {
    log.Fatal("Could not create TLS certificate.")
  }

  docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)
  if err != nil {
    log.Fatal("Error initializing docker: ", err)
  }

  fmt.Println("Starting!")
  c := &Containers{
    Client: *docker,
    Hostname: "whoa-daddy",
    Image: "ipython/scipystack",
    Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
  }
  c.Start()
  fmt.Println(c.ContainerId)
}
