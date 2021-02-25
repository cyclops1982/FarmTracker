package main

import (
	"encoding/binary"
	"encoding/json"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/cyclops1982/farmtracker/enhancedconn"
	"github.com/cyclops1982/farmtracker/protobufs"
	"github.com/cyclops1982/farmtracker/loramsgstructs"
	proto "github.com/golang/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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



func HandleClient(con enhancedconn.EnhancedConn) {
	defer con.Close()
	var msgLength uint16
	var err error

	msgLength = con.ReadLength();
	if msgLength == 0 {
		log.Printf("Client '%s' did not send correct length. Dropping.", con.RemoteAddr())
		return; 
	}

	initialdata := make([]byte, msgLength)
	_, err = con.ReadBytes(initialdata, 5)
	if err != nil {
		log.Printf("Client '%s' didn't send expected amount of data. Dropping.", con.RemoteAddr());
		return; 
	}

	msgReq := &protobufs.MessagesRequest{}

	err = proto.Unmarshal(initialdata, msgReq);
	if err != nil {
		log.Printf("Client '%s' didn't send expected protobuf message. Dropping.", con.RemoteAddr())
		return;
	}

	var filter string
	switch msgReq.DataToGet {
		case protobufs.MessagesRequest_LoraStatusV1:
			filter = "_status_"
		case protobufs.MessagesRequest_LoraUpdatesV1:
			filter = "_up_"
		case protobufs.MessagesRequest_LoraJoinV1:
			filter = "_join_"
		case protobufs.MessagesRequest_All:
			filter = ""
		default:
			filter = ""
	}

	log.Printf("New client '%s' requesting items with filter '%s' from %v\n", con.RemoteAddr(), filter, msgReq.DataSince.AsTime())

	// Now let's find some files in a seperate thread.
	ch := make(chan string)
	go FindFiles(ch, inputdir, filter, msgReq.DataSince.AsTime())

	// Read from the channel and send out the content.
	for file := range ch {

		//TODO: This needs some pattern to handle different file types based on their name.
		//Currently we only support the JSON format we receive from Chirpstack and assume it's the update message. (filter="_up_")

		filebytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("File '%s': Read error: %s\n", file, err)
			continue
		}
		// Verify that it's valid JSON.
		var jsonData interface{}
		err = json.Unmarshal(filebytes, &jsonData)
		if err != nil {
			log.Printf("File '%s': Failed to parse JSON. Error: %s\n", file, err)
			continue
		}

		// get the properties that we'd like to have.
		realData := jsonData.(map[string]interface{})
		devEUI, ok := realData["devEUI"].(string)
		if ok == false {
			log.Printf("File '%s': Failed to get devEUI.\n", file)
			continue
		}

		base64data, ok := realData["data"].(string)
		if ok == false {
			log.Printf("File '%s': Failed to get data.\n", file)
			continue
		}

		// convert the base64 string to a []byte
		bs, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			log.Printf("File '%s': Failed to decode base64 string '%s'.\n", file, base64data)
			continue
		}

		var loraMsg loramsgs.SodaqUniversalTracker
		byteReader := bytes.NewReader(bs)
		err = binary.Read(byteReader, binary.LittleEndian, &loraMsg)
		if err != nil {
			log.Printf("File '%s': Failed to unpack binary array from base64 data ('%s') into LoraMSG Structure\n", file, base64data);
			continue
		}

		msg := &protobufs.DeviceUpdate{}
		msg.DeviceIdentifier = &protobufs.DeviceIdentifier{}
		msg.DeviceIdentifier.Type = protobufs.DeviceIdentifier_DevEUI
		msg.DeviceIdentifier.Identifier = devEUI
		msg.Updated = timestamppb.New(time.Unix(int64(loraMsg.Unixtime), 0))
		msg.GPSCoordinates = &protobufs.Location{}
		msg.GPSCoordinates.Longitude = float64(loraMsg.Longitude)/10000000
		msg.GPSCoordinates.Latitude = float64(loraMsg.Latitude)/10000000
		msg.GPSCoordinates.Accuracy = 0;
		msg.BatteryVoltage = float32(((float32(loraMsg.RawVoltage)*10) + 3000)/1000 )
		msg.RawVoltage = uint32(loraMsg.RawVoltage)
		
		con.SendProtobufMsg(msg)

		log.Printf("File '%s' send to '%s'\n", file, con.RemoteAddr())
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
		go HandleClient(enhancedconn.EnhancedConn{con})
	}

}
