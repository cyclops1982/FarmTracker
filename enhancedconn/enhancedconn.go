package enhancedconn

import (
	"net"
	"log"
	"time"
	"encoding/binary"
	proto "github.com/golang/protobuf/proto"
)

type EnhancedConn struct {
	net.Conn
}

func (con *EnhancedConn) ReadBytes(allBytes []byte, timeout time.Duration) (int, error) {
	allReadBytes:=0

	if timeout != 0 {
		con.SetReadDeadline(time.Now().Add(timeout * time.Second))
	}	

	for {
		readBytes, err := con.Read(allBytes[allReadBytes:])
		allReadBytes += readBytes
		if err != nil {
			log.Printf("Failed to read %d bytes: %v\n", len(allBytes), err)
			return readBytes, err
		}
		if allReadBytes == len(allBytes) {
			break
		}
	}
	if timeout != 0 {
		con.SetReadDeadline(time.Time{})
	}
	return allReadBytes, nil
}


func (con *EnhancedConn) ReadLength() uint16 {
	bytes := make([]byte, 2)
	length, err := con.ReadBytes(bytes, 0)
	if (err != nil) {
		log.Printf("Failed to read length: %v\n", err);
		return 0
	}
	if length != len(bytes) {
		log.Printf("Didn't read %d bytes, only got %d: %v\n", len(bytes), length, err)
		return 0
	}

	return binary.BigEndian.Uint16(bytes)
}


func (con *EnhancedConn) SendProtobufMsg(msg proto.Message) {
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Failed to Marshal protobuf message. Error: %v\n", err)
	}

	// Send it over the network
	lengthMsg := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthMsg, uint16(len(out)))
	con.Write(lengthMsg)
	con.Write(out)
}