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
)

func Handle404(w http.ResponseWriter, r *http.Request) {
	log.Println("404 request")
	filedata, err := ioutil.ReadFile("404.html")
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


	// Format the filename. Time & event.
	filename := fmt.Sprintf("%s_%s.json", time.Now().Format(time.RFC3339Nano), event)

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
	return
}


func main() {
	log.Println("Starting http service....")

	router := mux.NewRouter()
	router.Path("/").Queries("event", "{event}").HandlerFunc(HandleRequest)
	router.HandleFunc("/", Handle404)

	srv := &http.Server{
		Handler: router,
		Addr: ":8989",
		// Timeouts to avoid overloading the server
		WriteTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe());
}
