package main

import (
  "github.com/odewahn/swarm-manager/manager"
  "github.com/odewahn/swarm-manager/models"
  "github.com/odewahn/swarm-manager/db"
  "fmt"
  "net/http"
  "os"
  "github.com/gorilla/mux"

)


// curl -X POST 127.0.0.1:8000/spawn
func Spawn(w http.ResponseWriter, r *http.Request) {
  // Make sure we can only be called with an HTTP POST request.
  /*
  if r.Method != "POST" {
    w.Header().Set("Allow", "POST")
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }
  */

  m := &models.Container{
    Image: "ipython/scipystack",
    Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
  }

  status := make(chan string)
  go manager.Start(m, status)

  <-status //block until we get a message back that the status record is ready

  fmt.Fprintf(w, m.Serialize())

}



func ListContainer(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  hostname := vars["hostname"]
  c := db.GetContainer(hostname)
  fmt.Fprintln(w, c.Serialize())
}



func KillContainer(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  hostname := vars["hostname"]
  c := db.GetContainer(hostname)
  status := make(chan string)
  go manager.Kill(&c, status)
  <-status  // block until the status updates
  fmt.Fprintln(w, c.Serialize())
}
