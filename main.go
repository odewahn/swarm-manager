package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/odewahn/swarm-manager/db"
	"github.com/odewahn/swarm-manager/manager"
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
func cors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		h.ServeHTTP(w, r)
	}
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.Header().Set("Access-Control-Allow-Origin", "*")

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
	// mux.Headers("Content-Type", "application/json")
	// mux.Headers("Access-Control-Allow-Headers", "origin, content-type, accept")
	mux.HandleFunc("/api/spawn/", use(Spawn, cors)).Methods("POST", "OPTIONS")
	http.Handle("/api/spawn/", mux)

	mux.HandleFunc("/api/container/{hostname}", use(ListContainer, cors)).Methods("GET")
	http.Handle("/api/container/{hostname}", mux)

	mux.HandleFunc("/api/container/{hostname}/kill", use(KillContainer, cors)).Methods("GET")
	http.Handle("/api/container/{hostname}/kill", mux)

	mux.HandleFunc("/api/containers", use(ListContainers, cors)).Methods("GET")
	http.Handle("/api/containers", mux)

	mux.HandleFunc("/manage", use(ManageContainers, cors)).Methods("GET")
	http.Handle("/manage", mux)

	mux.HandleFunc("/api/stats", use(Stats, cors)).Methods("GET")
	http.Handle("/api/stats", mux)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, mux); err != nil {
		fmt.Println(err.Error())
	}

}
