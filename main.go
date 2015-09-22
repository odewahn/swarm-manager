package main

import (
  "github.com/odewahn/swarm-manager/manager"
  "github.com/joho/godotenv"
  "fmt"
  "time"
  "os"
  "log"
)

func main() {

  // Load the environment variables we need
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  manager.Init()

  m := &manager.Container{
    Image: "ipython/scipystack",
    Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
    ContainerId: "suspicious_emu_5",
  }

  go m.Kill()

  for {
    fmt.Print(".")
    time.Sleep(500 * time.Millisecond)
  }
}
