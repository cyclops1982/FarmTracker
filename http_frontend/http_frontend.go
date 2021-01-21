package main

import (
	"os"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"fmt"
	"flag"
	"path/filepath"
	"html/template"
	"encoding/json"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

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

	// If we don't have a file path, we default to /index.html - the FileServer would do this too, but because we check
	// if the file exists below, it means that actualy setting the path to a *real* file makes the code below easier.
	if path == "/" {
		path = "/index.html";
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)
	
	// check whether a file exists at the given path
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(path)
	
	if os.IsNotExist(err) || fileInfo.IsDir() {
		// If we can't find the file, then we serve the 404 page, which is also template based.
		w.WriteHeader(http.StatusNotFound);
		fileInfo, _ = os.Stat(filepath.Join(h.staticPath, "404.html"))
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if the file we need to serve ends on .html, then it must be a template and we parse that below.
	filename := fileInfo.Name()
	if filepath.Ext(filename) == ".html" {
		var parsedFile *template.Template
		globPath := filepath.Join(h.staticPath, "*.html")
		parsedFile, err := template.ParseGlob(globPath);
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = parsedFile.ExecuteTemplate(w, filename[:len(filename) - 5] ,nil) // -5 == len(".html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	 } else {
		// otherwise, use http.FileServer to serve the static content like CSS/JS stuff
		http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
		return
	}
}


type BatReading struct {
	Date time.Time
	Voltage float32
}


func CreateConnection(dsn *string) *sql.DB {
	var pool *sql.DB
	var err error
	// Setup SQL connection
	pool, err = sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to SQL: %v\n", err)
	}
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)

	return pool
}

func CheckConnection(ctx context.Context, pool *sql.DB) {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("PingContext() failed - Unable to connect to database: %v\n", err)
	}

	if err := pool.Ping(); err != nil {
		log.Fatalf("Ping() failed - unable to connect to databsae: %v\n", err)
	}
}

func HandleJSONRequest(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	pool := CreateConnection(&sqlConString)

	var stuff []BatReading

	rows, err := pool.Query("SELECT LoggedOn, (((RawValue*10) + 3000)/1000 ) AS RealVolt FROM BatteryStatus ORDER BY LoggedOn")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError);
		log.Fatal(err);

		return;
	}
	defer rows.Close()
	for rows.Next() {
	var LoggedOn time.Time
	var Voltage float32
	err := rows.Scan(&LoggedOn, &Voltage);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError);
		log.Fatal(err);
		return;
	}
	stuff = append(stuff, BatReading{LoggedOn, Voltage})
}

	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(stuff)
}


var sqlConString string

func main() {
	// parameters
	var contentDir = flag.String("contentdir", "html", "Directory with static content to host.")
	var httpPort = flag.Int("port", 9191, "The port to bind on for the HTTP server.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind on.")
	flag.StringVar(&sqlConString, "sqlconstring", "farmtracker:MyGreatPassword@tcp(localhost)/FarmTracker?parseTime=true", "The DSN Connection String to use to connect to the MySQL DB.")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/api/json/{what}.json", HandleJSONRequest);
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
