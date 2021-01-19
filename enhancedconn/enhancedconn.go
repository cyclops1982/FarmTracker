package enhancedconn

import (
	"net"
)

type EnhancedConn struct {
	net.Conn
}

func (con *EnhancedConn) ReadBytes(allBytes []byte) (int, error) {
	allReadBytes:=0

	for {
		readBytes, err := con.Read(allBytes[allReadBytes:])
		allReadBytes += readBytes
		if err != nil {
			return readBytes, err
		}
		if allReadBytes == len(allBytes) {
			break
		}
	}
	return allReadBytes, nil
}
