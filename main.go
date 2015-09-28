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
  "encoding/base64"
  "strings"
  "os"
)

var (
  HTTPAddr = flag.String("http", "0.0.0.0:3000", "Address to listen for HTTP requests on")
)


// The following 2 functions are used to do basic auth and come from
//   https://gist.github.com/elithrar/9146306
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != os.Getenv("USERNAME") {
			http.Error(w, "Not authorized", 401)
			return
		}

    if pair[1] != os.Getenv("PASSWORD") {
      http.Error(w, "Not authorized", 401)
      return
    }

		h.ServeHTTP(w, r)
	}
}



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

  mux.HandleFunc("/spawn", use(Spawn, basicAuth)).Methods("POST")
  http.Handle("/spawn", mux)

  mux.HandleFunc("/container/{hostname}", use(ListContainer, basicAuth)).Methods("GET")
  http.Handle("/container/{hostname}", mux)


  mux.HandleFunc("/container/{hostname}/kill", use(KillContainer, basicAuth)).Methods("GET")
  http.Handle("/container/{hostname}/kill", mux)


  mux.HandleFunc("/containers", use(ListContainers, basicAuth)).Methods("GET")
  http.Handle("/containers", mux)


  mux.HandleFunc("/manage", use(ManageContainers, basicAuth)).Methods("GET")
  http.Handle("/manage", mux)

  // Start the HTTP server!
   fmt.Println("HTTP server listening on", *HTTPAddr)
   if err := http.ListenAndServe(*HTTPAddr, mux); err != nil {
     fmt.Println(err.Error())
   }


}
