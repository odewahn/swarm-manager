package main

import (
  "github.com/odewahn/swarm-manager/manager"
  "github.com/odewahn/swarm-manager/db"
  "github.com/joho/godotenv"
  "log"
  "flag"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
)

var (
  HTTPAddr = flag.String("http", "0.0.0.0:3000", "Address to listen for HTTP requests on")
)

func main() {

  // Load the environment variables we need
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  flag.Parse()

  manager.Init()
  db.Init()

  mux := mux.NewRouter()

  mux.HandleFunc("/spawn", Spawn).Methods("POST")
  mux.HandleFunc("/container/{hostname}", ListContainer).Methods("GET")
  mux.HandleFunc("/container/{hostname}/kill", KillContainer).Methods("GET")
  mux.HandleFunc("/containers", ManageContainers).Methods("GET")

  // Start the HTTP server!
   fmt.Println("HTTP server listening on", *HTTPAddr)
   if err := http.ListenAndServe(*HTTPAddr, mux); err != nil {
     fmt.Println(err.Error())
   }


}
