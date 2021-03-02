package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"github.com/google/uuid"
	"strings"
	"flag"
)

func GenUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	log.Println("404 request")
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

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event, eventexists := vars["event"]
	if eventexists == false {
		log.Fatalln("Request should have a querystring parameter for 'event', but we couldn't find it.")
		return
	}
	if event == "" {
		w.WriteHeader(400)
		w.Write([]byte("We're expecting a value in the 'event' parameter."))
		log.Println("400 - Event is empty.")
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("This web service expects a POST."))
		log.Println("405 - Method wasn't POST")
		return
	}


	dirname := fmt.Sprintf("%s%s", g_outputDir, time.Now().Format("2006/01/02/"))
	filename := fmt.Sprintf("%s%s_%s_%s.json", dirname, time.Now().Format(time.RFC3339Nano), event, GenUUID())

	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		log.Println("500 - ERROR - Failed to create directory ", dirname, "Error:", err)
		w.WriteHeader(500)
		w.Write([]byte("Couldn't create directory for storage."))
		return
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		log.Println("500 - ERROR - Failed to create file: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Couldn't create the local file, please try again."))
		return
	}
	defer f.Close()

	copiedBytes, copyErr := io.Copy(f, r.Body)
	if copyErr != nil {
		log.Println("500 - ERROR - Failed to copy request body: ", copyErr)
		w.WriteHeader(500)
		w.Write([]byte("Couldn't copy request body."))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(string(copiedBytes)))
	log.Println("200 - Received", copiedBytes," bytes. Writen to ", filename)
	return
}

var g_outputDir string
func main() {
	// parameters
	var httpPort = flag.Int("port", 9090, "The port to bind on for the HTTP server. Default is 9090.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind on, default is 0.0.0.0.")
	flag.StringVar(&g_outputDir, "outputdir", "dumps/", "The path to store files. A directory structure YYYY/MM/DD/ will be created in this folder.")
	flag.Parse()

	// Check base path
	err := os.MkdirAll(g_outputDir, 0755)
	if err != nil {
		log.Fatal("Couldn't create outputdir.", err)
	}

	router := mux.NewRouter()
	router.Path("/").Queries("event", "{event}").HandlerFunc(HandleRequest)
	router.HandleFunc("/", Handle404)

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
