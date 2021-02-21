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
	"encoding/binary"
	"io/ioutil"
	"bytes"
	"strings"
)

func FindFiles(ch chan string, path string, filter string, unixtime time.Time) {
	foundFiles := make(map[string]bool)
	for {
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if info.IsDir() == false && (filter == "" || strings.Contains(info.Name(), filter))  {
				underscoreloc := strings.IndexRune(info.Name(), '_')
				if underscoreloc == -1 {
					log.Printf("Couldn't get date/time from file: %s\n", p)
					return nil
				}
				filedatetime, err := time.Parse(time.RFC3339Nano, info.Name()[:underscoreloc])
				if (err != nil) {
					log.Printf("Couldn't parse date/time from file: %s\n", p)
					return nil
				}
				if (filedatetime.After(unixtime) && foundFiles[p] == false) {
					ch <- p
					foundFiles[p] = true
				}
			}
			return nil
		})
		time.Sleep(1 * time.Second)
	}
}


func HandleClient(con net.Conn) {
	defer con.Close()

	var err error

	//TODO: Move this timestamp and what the client wants to receive into a protobuf. That will allow us to enhance this protocol later. We simply don't know what we'll do later.
	// We're reading a few bytes to get some data that we need.
	// That would first be a uint64 (for a unixtime)
	tmpUTime := make([]byte, 8)
	_, err = con.Read(tmpUTime)
	if err != nil {
		log.Printf("Expected 8 bytes (uint64) from %s. Got an error: %v. Disconnecting.\n", con.RemoteAddr(), err)
		return
	}
	timeInt64 := int64(binary.LittleEndian.Uint64(tmpUTime))
	unixTime := time.Unix(timeInt64, 0)
	log.Printf("The time: %v - %v\n", timeInt64, unixTime)

	tmp := make([]byte, 30)
	_, err = con.Read(tmp)
	if err != nil {
		log.Println("Failed to read from client",con.RemoteAddr(),". Disconnecting. Error: ", err)
		return
	}
	// check if we have a dot, as that indicates the end of our command.
	if bytes.ContainsRune(tmp, '.') == false {
		con.Write([]byte("Sorry, i didn't get that. Bye."))
		return
	}

	// Remove all chars, so we compare to a nice string.
	replacer := strings.NewReplacer("\r", "", "\n", "", ".", "", "\x00", "")
	filterinput := replacer.Replace(string(tmp[:]))

	var filter string
	switch filterinput {
		case "up":
			filter = "_up_"
		case "status":
			filter = "_status_"
		case "join":
			filter = "_join_"
		default:
			filter = ""
	}
	if filter != "" {
		log.Println("Filtering returned items on ", filter)
	}

	// Now let's find some files in a seperate thread.
	ch := make(chan string)
	go FindFiles(ch, inputdir, filter, unixTime)

	// Read from the channel and send out the content.
	for file := range ch {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("Failed to read ", file, ". Error:", err)
			continue
		}
		// Verify that it's valid JSON.
		var jsonData interface{}
		err = json.Unmarshal(bytes, &jsonData)
		if err != nil {
			log.Printf("Failed to parse JSON. Ignoring file '%s'. Error: %s\n", file, err)
			continue
		}
		/* 
		Example of filtering:
		realData := jsonData.(map[string]interface{})
		if realData["applicationID"] != "1" {
			log.Println("JSON payload is not for application 1. Skipping.")
			continue
		}*/
		// Send out 2 bytes that tell us how long the message will be
		length := make([]byte, 2)
		binary.BigEndian.PutUint16(length, uint16(len(bytes)))
		_, err = con.Write(length)
		writenbytes, err := con.Write(bytes)
		if err != nil {
			log.Println("Failed to write bytes. Disconnecting. Error was: ", err)
			return
		} else {
			log.Printf("Wrote %d bytes to %s - %s.\n", writenbytes, con.RemoteAddr(), file)
		}
	}
}


var inputdir string
func main() {
	flag.StringVar(&inputdir, "inputdir", "dumps/", "The directory to process.")
	var tcpPort = flag.Int("port", 29000, "The port to listen on. Default is 29000.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind to, default is 0.0.0.0.")
	flag.Parse()

	addr:=fmt.Sprintf("%s:%d", *ipAddress, *tcpPort)
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Couldn't open TCP socket.", err)
	}

	// Wait for clients to connect
	log.Println("Listening on ", addr)
	for {
		con, err := sock.Accept()
		if err != nil {
			log.Println("Failed to accept: ", err)
			continue;
		}
		log.Println("Client ",con.RemoteAddr()," connected. Starting Thread.")
		go HandleClient(con)
	}

}
