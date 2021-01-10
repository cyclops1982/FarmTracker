package main

import (
	"os"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
	"flag"
	"path/filepath"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("404 - %s request for '%s' from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	filedata, err := ioutil.ReadFile("html/404.html")
	if err != nil {
		log.Println("ERROR - Couldn't read 404.html file.. returning standard http.NotFound")
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(404)
	w.Write(filedata)
	return
}

type PageHandler struct {
	staticPath string
}

func (h PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		NotFoundHandler(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}


var g_outputDir string
func main() {
	// parameters
	var contentDir = flag.String("contentdir", "html", "Directory with static content to host.")
	var httpPort = flag.Int("port", 9191, "The port to bind on for the HTTP server.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind on.")
	flag.Parse()

	router := mux.NewRouter()
	
	router.PathPrefix("/").Handler(PageHandler{staticPath: *contentDir})

	log.Printf("Starting HTTP server on %s:%d\n", *ipAddress, *httpPort)
	srv := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf("%s:%d", *ipAddress, *httpPort),
		// Timeouts to avoid overloading the server
		WriteTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe());
}
