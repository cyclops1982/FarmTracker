package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cyclops1982/farmtracker/enhancedconn"
	"github.com/cyclops1982/farmtracker/protobufs"
	proto "github.com/golang/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	var err error
	var ipAddress = flag.String("server", "127.0.0.1", "The IP address (or hostname) of the server to connect to.")
	var tcpPort = flag.Int("port", 29000, "The port to use for the server connection.")
	//var sqlConString  = flag.String("sqlconstring", "farmtracker:MyGreatPassword@tcp(localhost)/FarmTracker", "The DSN Connection String to use to connect to the MySQL DB.")
	//TODO: Add a 'from' parameter that just takes an amount of hours
	var fromUnixtime = flag.Int64("fromUnixtime", 0, "Set the unix timestamp from which we should receive messages.")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *ipAddress, *tcpPort)
	
	tmpcon, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to %s. Exiting.", addr)
	}
	con := enhancedconn.EnhancedConn{tmpcon}

	req := &protobufs.MessagesRequest{}
	req.DataToGet = protobufs.MessagesRequest_LoraUpdatesV1
	req.DataSince = timestamppb.New(time.Unix(*fromUnixtime, 0))

	out, err := proto.Marshal(req)
	if err != nil {
		log.Fatal("wUT", err)
	}

	lengthMsg := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthMsg, uint16(len(out)))
	con.Write(lengthMsg)
	con.Write(out)
	
	for {
		length, err := con.ReadLength();
		if length == 0 || err != nil {
			log.Printf("Server didn't send us a message length.\n")
			return; 
		}

		initialdata := make([]byte, length)
		_, err = con.ReadBytes(initialdata, 5)
		if err != nil {
			log.Printf("Didn't receive expected amount of data.\n");
			return; 
		}


		msg := &protobufs.DeviceUpdate{}
		err = proto.Unmarshal(initialdata, msg)
		log.Printf("Message: %v\n", msg)

	}
}
