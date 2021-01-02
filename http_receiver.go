package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
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
	log.Println("Normal request")
	http.NotFound(w, r)
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
