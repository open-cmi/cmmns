package transport

import (
	"net"
	"time"
)

// Header tlv header
type Header struct {
	Type uint16
	Len  uint16
}

// TCPRequest request msg
func TCPRequest(msg []byte) (conn *net.TCPConn, err error) {

	rsvr, err := net.ResolveTCPAddr("tcp", gConf.TCPServer)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP("tcp", nil, rsvr)
	if err != nil {
		return conn, err
	}
	now := time.Now()

	conn.SetDeadline(now.Add(10 * time.Second))
	conn.SetReadDeadline(now.Add(10 * time.Second))
	conn.SetWriteDeadline(now.Add(10 * time.Second))
	_, err = conn.Write(msg)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
