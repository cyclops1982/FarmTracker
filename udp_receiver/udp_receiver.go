package main

import (
	"time"
	"flag"
	"net"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/google/uuid"
//	"github.com/cyclops1982/farmtracker/enhancedconn"
	
)

func GenUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

func ReceiveUDP(con net.Conn ) {
	client := con.RemoteAddr()
	dirname := fmt.Sprintf("%s%s", g_outputDir, time.Now().Format("2006/01/02/"))
	filename := fmt.Sprintf("%s%s%s_%s_%s.raw", dirname, client.Network(), time.Now().Format(time.RFC3339Nano), strings.Split(client.String(),":")[0], GenUUID())

	data := make([]byte, 100)
	readBytes, err := con.Read(data)
	log.Printf("Read %d bytes. Error: %v\n", readBytes, err)
	log.Printf("Writing to ", filename)
	log.Println(data)

}


var g_outputDir string
func main() {
	// parameters
	var udpPort = flag.Int("port", 7084, "The port to bind to. Default is 7084.")
	var ipAddress = flag.String("address", "0.0.0.0", "The IP address to bind on, default is 0.0.0.0.")
	flag.StringVar(&g_outputDir, "outputdir", "dumps/", "The path to store files. A directory structure YYYY/MM/DD/ will be created in this folder.")
	flag.Parse()

	// Check base path
	err := os.MkdirAll(g_outputDir, 0755)
	if err != nil {
		log.Fatal("Couldn't create outputdir.", err)
	}

	
	addr:=fmt.Sprintf("%s:%d", *ipAddress, *udpPort)
	sock, err := net.ListenPacket("udp4", addr)
	if err != nil {
		log.Fatal("Couldn't open UDP socket.", err)
	}
	defer sock.Close()

	// Wait for clients to connect
	log.Println("Listening on ", addr)
	for {
		data := make([]byte, 100)
		readBytes, addr, err := sock.ReadFrom(data)
		log.Printf("Read %d bytes from %s. Error: %v\n", readBytes, addr, err)
		log.Println(data)

//		go ReceiveUDP(con)
	}

}
