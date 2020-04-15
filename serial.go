package main

import (
	"io"
	"time"
	"log"
	"github.com/pkg/errors"
	"github.com/jacobsa/go-serial/serial"
)

var verbose = false

func SerialInit(port string, v bool) (stream io.ReadWriteCloser, err error) {
	verbose = v
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
	if verbose {log.Printf(" --> serial.Open(%#v)\n",options)}
	stream,err = serial.Open(options);
	if err != nil {
		return nil, errors.Wrapf(err,"Could not open serial port")
	}
	return stream, nil
}
	
func SerialReadByte(stream io.ReadWriteCloser, timeoutMS uint, purpose string) (b byte, err error) {
	var arr []byte = make([]byte,1);
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		if verbose {log.Printf("       serial.Read (%d ms) - %s\n",timeoutMS,purpose)}
		_,readerr := stream.Read(arr); 
		if verbose {log.Printf(" %2x <-- serial.Read (%d ms)\n",arr[0],timeoutMS)}
		c <- readerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			return 0,errors.Wrap(err, "failed to read byte")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		if verbose {log.Printf("   read TIMEOUT at %d ms\n", timeoutMS)}
		return 0,errors.Errorf("TIMEOUT: read timed out at %d ms", timeoutMS)
	}
	return arr[0], nil
}

func SerialWriteByte(stream io.ReadWriteCloser, timeoutMS uint, b byte, purpose string) (err error) {
	var arr []byte = make([]byte,1);
	arr[0] = b;
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		if verbose {log.Printf(" --> %2x serial.Write (%d ms) - %s\n",arr[0], timeoutMS,purpose)}
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
		if verbose {log.Printf("   write TIMEOUT at %d ms\n", timeoutMS)}
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}

func SerialWriteBytes(stream io.ReadWriteCloser, timeoutMS uint, arr []byte, purpose string) (err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		if verbose {log.Printf(" --> %2x serial.Write (%d ms) - %s\n",arr,timeoutMS,purpose)}
		_,writeerr := stream.Write(arr); 
		c <- writeerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			return errors.Wrap(err, "failed to write bytes")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		if verbose {log.Printf("   write TIMEOUT at %d ms\n", timeoutMS)}
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}
