package main

import (
	"log"
	"flag"
	"path/filepath"
	"os"
	"time"
	"net"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func FindFiles(ch chan string, path string) {
	foundDirs := make(map[string]bool)
	for {
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if info.IsDir() == false {
				if foundDirs[p] == false {
					ch <- p
					foundDirs[p] = true
				}
			}
			return nil
		})
		time.Sleep(1 * time.Second)
	}
}


func HandleClient(con net.Conn) {
	defer con.Close()
	ch := make(chan string)
	go FindFiles(ch, inputdir)

	for file := range ch {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("Failed to read ", file, ". Error:", err)
			continue
		}
		// Verify that it's valid JSON.
		// Possibly we'd do some filtering too.
		var jsonData interface{}
		err = json.Unmarshal(bytes, &jsonData)
		if err != nil {
			log.Println("Failed to parse JSON in ",file,". Error: ", err)
			continue
		}
		realData := jsonData.(map[string]interface{})
		if realData["applicationID"] != "1" {
			log.Println("JSON payload is not for application 1. Skipping.")
			continue
		}
		writenbytes, err := con.Write(bytes)
		if err != nil {
			log.Println("Failed to write bytes: ", err)
		} else {
			log.Println("Wrote ",writenbytes," to ", con.RemoteAddr())
		}
	}
}


var inputdir string
func main() {
	flag.StringVar(&inputdir, "inputdir", "dumps/", "The directory to process.")
	var tcpPort = flag.Int("port", 29000, "The port to listen on. Default is 29000.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind to, default is 0.0.0.0.")
	flag.Parse()

	sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ipAddress, *tcpPort))
	if err != nil {
		log.Fatal("Couldn't open TCP socket.", err)
	}

	// Wait for clients to connect
	for {
		con, err := sock.Accept()
		if err != nil {
			log.Println("Failed to accept: ", err)
			continue;
		}
		go HandleClient(con)
	}

}
