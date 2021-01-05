package main

import (
	"log"
	"net"
	"encoding/json"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
)



func main() {
	var ipAddress = flag.String("server", "127.0.0.1", "The IP address (or hostname) of the server to connect to. Default is 127.0.0.1.")
	var tcpPort = flag.Int("port", 29000, "The port to use for the server connection. Default is 29000.")
	//var sqlConString  = flag.String("sqlconstring", "FarmTracker:MyGreatPassword@localhost/FarmTracker", "The DSN Connection String to use to connect to the MySQL DB.")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *ipAddress, *tcpPort)
	con, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to %s. Exiting.", addr)
	}

	con.Write([]byte("up.")) // for now, just send a '.' so we get everything.

	msgLength := make([]byte, 2)
	for {
		nBytes, err := con.Read(msgLength)
		msgLengthUint16 := binary.BigEndian.Uint16(msgLength)
		log.Printf("Length of message that's coming: %d\n", msgLengthUint16)
		if nBytes != 2 {
			log.Fatal("We really expect 2 bytes for a messagelength.")
		}
		msgData := make([]byte, msgLengthUint16)
		nBytes, err = con.Read(msgData)
		if err != nil {
			log.Println("Failed to read:", err)
			continue
		}
		if (nBytes != int(msgLengthUint16)) {
			log.Fatal("Really expected the correct amount of bytes...")
		}

		// Chop the buffer to something smaller so we can correctly convert it to JSON
		var jsonData interface{}
		err = json.Unmarshal(msgData, &jsonData)
		if err != nil {
			log.Println("Failed to parse JSON:", err)
			continue
		}
		realData := jsonData.(map[string]interface{})
		base64data, ok := realData["data"].(string)
		if ok == false {
			log.Println("Failed to convert data to string.")
			continue
		}
		data, err := base64.StdEncoding.DecodeString(base64data)
		deveui := realData["devEUI"]
		log.Printf("devEUI: %s\nBase64data:%s\nDATA: %s", deveui, base64data, data)
	}
}
