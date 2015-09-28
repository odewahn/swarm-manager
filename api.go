package main

import (
  "github.com/odewahn/swarm-manager/manager"
  "github.com/odewahn/swarm-manager/models"
  "github.com/odewahn/swarm-manager/db"
  "fmt"
  "net/http"
  "os"
  "github.com/gorilla/mux"
  "html/template"
)

type SpawnRequest struct {
  Image string
}

func Spawn(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  image := r.FormValue("image")
  if len(image) == 0 {
    image = "ipython/scipystack"
  }

  user := r.FormValue("user")
  if len(user) == 0 {
    user = "odewahn"
  }

  m := &models.Container{
    Image: image,
    User: user,
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

func ManageContainers(w http.ResponseWriter, r *http.Request) {
  containers := db.GetContainers()
  t, _ := template.New("index").Parse(`
   <html>
      <table>
        {{ range .}}
         <tr>
            <td>
               <a target=_blank href="http://{{.Url}}">{{.Url}}</a>
            </td>
            <td>
               {{.Image}}
            </td>
            <td>
               {{.User}}
            </td>
            <td>
               {{.Status}}
            </td>
            <td>
               {{.StartTime}}
            </td>
            <td>
               <a href="/container/{{.Hostname}}/kill">Kill</a>
            </td>
          </tr>
         {{end}}
      </table>
  </html>
  `)
  t.Execute(w, containers)
}
