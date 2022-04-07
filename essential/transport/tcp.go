package transport

import (
	"net"
)

// Header tlv header
type Header struct {
	Type uint16
	Len  uint16
}

var TCPServer string = "localhost:30017"

// TCPRequest request msg
func TCPRequest(msg []byte) (conn *net.TCPConn, err error) {
	rsvr, err := net.ResolveTCPAddr("tcp", TCPServer)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP("tcp", nil, rsvr)
	if err != nil {
		return conn, err
	}

	_, err = conn.Write(msg)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
