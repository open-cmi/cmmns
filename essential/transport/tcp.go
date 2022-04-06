package transport

import (
	"net"
)

// Header tlv header
type Header struct {
	Type uint16
	Len  uint16
}

var XSNOSServer string = "192.168.1.104:30017" //"192.168.56.103:30017"

// TCPRequest request msg
func TCPRequest(msg []byte) (conn *net.TCPConn, err error) {
	rsvr, err := net.ResolveTCPAddr("tcp", XSNOSServer)
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
