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


func (con *EnhancedConn) ReadBytes(data []byte, timeout time.Duration) (int, error) {
	allReadBytes := 0
	if timeout != 0 {
		con.SetReadDeadline(time.Now().Add(timeout * time.Second))
		defer con.SetReadDeadline(time.Time{})
	}
	return io.ReadFull(data)
}


func (con *EnhancedConn) ReadLength() (uint16, error) {
	data := make([]byte, 2)
	length, err := io.ReadFull(data)
	if err != nil {
		return 0, err
	}
	if length != len(data) {
		log.Printf("Didn't read %d bytes, only got %d.\n", len(data), length)
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}


func (con *EnhancedConn) SendProtobufMsg(msg proto.Message) error {
	out, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	lengthMsg := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthMsg, uint16(len(out)))
	if _, err := con.Write(lengthMsg); err != nil {
		return err
	}
	_, err = con.Write(out)
	return err
}