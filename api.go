package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/odewahn/swarm-manager/db"
	"github.com/odewahn/swarm-manager/manager"
	"github.com/odewahn/swarm-manager/models"
)

type SpawnRequest struct {
	Image string
}

func Spawn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image := r.FormValue("image")
	if len(image) == 0 {
		image = "zischwartz/notebook"
	}

	user := r.FormValue("user")
	if len(user) == 0 {
		user = "odewahn"
	}

	m := &models.Container{
		Image:      image,
		User:       user,
		Domainname: os.Getenv("THEBE_SERVER_BASE_URL"),
	}

	status := make(chan string)
	go manager.Start(m, status)

	<-status //block until we get a message back that the status record is ready

	fmt.Fprintf(w, m.Serialize()+"\n")

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
	<-status // block until the status updates
	fmt.Fprintln(w, c.Serialize())
}

func ListContainers(w http.ResponseWriter, r *http.Request) {
	containers := db.GetContainers()
	out, err := json.MarshalIndent(containers, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(out)+"\n")
}

func Stats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}\n")
}

func ManageContainers(w http.ResponseWriter, r *http.Request) {
	containers := db.GetContainers()
	t, _ := template.New("index").Parse(`
   <html>
      <h1>Swarm Manager</h1>
      <h2>Launch a container</h2>
      <form action="/api/spawn/" method="POST">
         Image: <input type="text" name="image" value="ipython/scipystack"/><br>
         User: <input type="text" name="user" value="odewahn"/><br>
         <input type="submit"/>
      </form>
      <h2>Active Containers</h2>
      <table>
        {{ range .}}
				{{ if .IsActive }}
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
               <a href="/api/container/{{.Hostname}}/kill">Kill</a>
            </td>
          </tr>
					{{ end }}
         {{end}}
      </table>
  </html>
  `)
	t.Execute(w, containers)
}
