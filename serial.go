package main

import (
	"io"
	"time"
	"github.com/pkg/errors"
	"github.com/jacobsa/go-serial/serial"
)

func SerialInit(port string) (stream io.ReadWriteCloser, err error) {
	options := serial.OpenOptions{
		PortName: port,
		BaudRate: 9600,
		ParityMode: serial.PARITY_NONE,
		RTSCTSFlowControl: true,
		InterCharacterTimeout: 500,
		MinimumReadSize: 1,
		DataBits: 8,
		StopBits: 1,
	}
	stream,err = serial.Open(options);
	if err != nil {
		return nil, errors.Wrapf(err,"Could not open serial port")
	}
	return stream, nil
}
	
func SerialReadByte(stream io.ReadWriteCloser, timeoutMS uint) (b byte, err error) {
	var arr []byte = make([]byte,1);
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		_,readerr := stream.Read(arr); 
		c <- readerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			return 0,errors.Wrap(err, "failed to read byte")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		return 0,errors.Errorf("read timed out at %d ms", timeoutMS)
	}
	return arr[0], nil
}

func SerialWriteByte(stream io.ReadWriteCloser, timeoutMS uint, b byte) (err error) {
	var arr []byte = make([]byte,1);
	arr[0] = b;
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		_,writeerr := stream.Write(arr); 
		c <- writeerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			return errors.Wrap(err, "failed to write byte")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		return errors.Errorf("write timed out at %d ms", timeoutMS)
	}
	return nil
}
