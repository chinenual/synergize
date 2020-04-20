package main

import (
	"fmt"
	"io"
	"time"
	"log"
	"github.com/pkg/errors"
	"github.com/jacobsa/go-serial/serial"
)

var verbose = false

func serialInit(port string, v bool) (stream io.ReadWriteCloser, err error) {
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
	
func serialReadByte(stream io.ReadWriteCloser, timeoutMS uint, purpose string) (b byte, err error) {
	var arr []byte = make([]byte,1);
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if verbose {log.Printf("       serial.Read (%d ms) - %s\n",timeoutMS,purpose)}
	
	c := make(chan error, 1)
	go func() {
		_,readerr := stream.Read(arr); 
		c <- readerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			return arr[0],errors.Wrap(err, "failed to read byte")
		}
		if verbose {log.Printf(" %02x <-- serial.Read (%d ms)\n",arr[0],timeoutMS)}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		if verbose {log.Printf("   read TIMEOUT at %d ms (%x)\n", timeoutMS,arr[0])}
		return arr[0],errors.Errorf("TIMEOUT: read timed out at %d ms", timeoutMS)
	}
	return arr[0], nil
}


func serialReadBytes(stream io.ReadWriteCloser, timeoutMS uint, num_bytes uint16, purpose string) (bytes []byte, err error) {
	var arr []byte = make([]byte,num_bytes);


	for i:= uint16(0); i < num_bytes; i++ {
		arr[i],err = serialReadByte(stream, timeoutMS, fmt.Sprintf("%s: %d",purpose,i))
		if err != nil {
			bytes = arr[0:i]
			err = errors.Wrap(err, "failed to read all bytes")
		}
	}
	bytes = arr
	return
}
	
func bugserialReadBytes(stream io.ReadWriteCloser, timeoutMS uint, num_bytes uint16, purpose string) (bytes []byte, err error) {
	var arr []byte = make([]byte,num_bytes);
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if verbose {log.Printf("       serial.ReadBytes (%d ms) %d - %s\n",timeoutMS,num_bytes,purpose)}

	var num_read int
	c := make(chan error, 1)
	go func() {
		var readerr error
		num_read, readerr = stream.Read(arr); 
		c <- readerr
	} ()
	
	select {
	case err := <-c:
		if err != nil {
			// truncate the array in case fewer than requested were read
			arr = arr[0:num_read]
			return arr,errors.Wrap(err, "failed to read byte")
		}
		if verbose {log.Printf(" %02x <-- serial.Read %d (%d ms)\n",arr,num_read,timeoutMS)}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		if verbose {log.Printf("   read TIMEOUT at %d ms (%x)\n", timeoutMS,arr)}
		return arr,errors.Errorf("TIMEOUT: read timed out at %d ms", timeoutMS)
	}
	return arr, nil
}

func serialWriteByte(stream io.ReadWriteCloser, timeoutMS uint, b byte, purpose string) (err error) {
	var arr []byte = make([]byte,1);
	arr[0] = b;
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if verbose {log.Printf(" --> %02x serial.Write (%d ms) - %s\n",arr[0], timeoutMS,purpose)}

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
		if verbose {log.Printf("   write TIMEOUT at %d ms\n", timeoutMS)}
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}

func serialWriteBytes(stream io.ReadWriteCloser, timeoutMS uint, arr []byte, purpose string) (err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if verbose {log.Printf(" --> %02x serial.WriteBytes (%d ms) - %s\n",arr,timeoutMS,purpose)}

	c := make(chan error, 1)
	go func() {
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
