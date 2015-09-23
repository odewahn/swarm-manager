package main

import (
  "github.com/odewahn/swarm-manager/manager"
  "github.com/joho/godotenv"
  "os"
  "log"
  "flag"
  "fmt"
  "time"
)

var (
  Action = flag.String("action", "START", "Action (START | KILL | NOOP)")
  Hostname = flag.String("hostname", "", "Hostname to kill")
)

func main() {

  // Load the environment variables we need
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  flag.Parse()

  manager.Init()


  m := &manager.Container{
    Image: "ipython/scipystack",
    Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
  }

  if *Action == "START" {
    go m.Start()
  }

  if *Action == "NOOP" {
    fmt.Println("doing a noop")
    go m.NoOp()
  }

  if *Action == "KILL" {
    m.Hostname = *Hostname
    go m.Kill()
  }

  for {
    fmt.Print(".")
    time.Sleep(500*time.Millisecond)
  }

}
