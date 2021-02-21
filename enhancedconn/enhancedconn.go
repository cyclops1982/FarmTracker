package enhancedconn

import (
	"net"
	"log"
	"time"
	"encoding/binary"
)

type EnhancedConn struct {
	net.Conn
}

func (con *EnhancedConn) ReadBytes(allBytes []byte) (int, error) {
	allReadBytes:=0

	con.SetReadDeadline(time.Now().Add(5 * time.Second))

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
	return allReadBytes, nil
}


func (con *EnhancedConn) ReadLength() uint16 {
	bytes := make([]byte, 2)
	length, err := con.ReadBytes(bytes)
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