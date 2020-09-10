package io

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"time"
)

type SocketIo struct {
	conn *net.TCPConn
}

func SocketInit(addr string) (s SocketIo, err error) {
	log.Printf(" --> socket.Open(%s)\n", addr)
	var tcpAddr *net.TCPAddr
	if tcpAddr, err = net.ResolveTCPAddr("tcp4", addr); err != nil {
		return
	}
	if s.conn, err = net.DialTCP("tcp", nil, tcpAddr); err != nil {
		return
	}
	return
}

func (s SocketIo) close() (err error) {
	if err = s.conn.Close(); err != nil {
		return
	}
	return
}

func (s SocketIo) readByte(timeoutMS uint) (b byte, err error) {
	if err = s.conn.SetDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
		return
	}
	var arr []byte = make([]byte, 1)
	var n int
	n, err = s.conn.Read(arr)
	if err != nil {
		return
	}
	if n != 1 {
		err = errors.New("TIMEOUT reading byte")
		return
	}
	b = arr[0]
	return
}

func (s SocketIo) readBytes(timeoutMS uint, num_bytes uint16) (bytes []byte, err error) {
	bytes = make([]byte, num_bytes)

	if err = s.conn.SetDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
		return
	}
	var n int
	n, err = s.conn.Read(bytes)
	if err != nil {
		return
	}
	if n != int(num_bytes) {
		err = errors.New("TIMEOUT reading bytes")
		return
	}
	return
}

func (s SocketIo) writeByte(timeoutMS uint, b byte) (err error) {
	if err = s.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
		return
	}
	var arr []byte = make([]byte, 1)
	var n int
	arr[0] = b
	n, err = s.conn.Write(arr)
	if err != nil {
		return
	}
	if n != 1 {
		err = errors.New("TIMEOUT writing byte")
		return
	}
	return
}

func (s SocketIo) writeBytes(timeoutMS uint, arr []byte) (err error) {
	if err = s.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
		return
	}
	var n int
	n, err = s.conn.Write(arr)
	if err != nil {
		return
	}
	if n != len(arr) {
		err = errors.New("TIMEOUT writing bytes")
		return
	}
	return
}
