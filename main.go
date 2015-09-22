package main

import (
  "fmt"
  "github.com/joho/godotenv"
  "github.com/samalba/dockerclient"
  "log"
  "os"
  "github.com/odewahn/swarm-manager/container"
  "time"
)

func main() {

  // Load the environment variables we need
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  tlsConfig, err := container.GetTLSConfig(os.Getenv("SWARM_CREDS_DIR"))
  if err != nil {
    log.Fatal("Could not create TLS certificate.")
  }

  docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)
  if err != nil {
    log.Fatal("Error initializing docker: ", err)
  }

  fmt.Println("Starting ", docker)

  //container.Init(docker)

  c := &container.Container{
    Hostname: "whoa-daddy",
    Image: "ipython/scipystack",
    Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
    ContainerId: "6e1111899edafd3b5d50486e59317ede3ea21da272c978ec834251fdb21a010b",
  }
  go c.Start()
  var input string
  for {
    fmt.Print(".")
    time.Sleep(500 * time.Millisecond)
  }
  fmt.Scanln(&input)
  fmt.Println(c.ContainerId)
}
