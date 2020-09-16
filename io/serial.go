package io

import (
	"io"
	"time"

	"github.com/chinenual/synergize/logger"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
)

type serialReadResponse struct {
	data byte
	err  error
}

type SerialIo struct {
	stream io.ReadWriteCloser
}

var readerChannel chan serialReadResponse

func SerialInit(port string, baudRate uint) (s SerialIo, err error) {
	options := serial.OpenOptions{
		PortName:              port,
		BaudRate:              baudRate,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     true,
		InterCharacterTimeout: 500,
		MinimumReadSize:       1,
		DataBits:              8,
		StopBits:              1,
	}
	logger.Infof(" --> serial.Open(%#v)\n", options)
	if s.stream, err = serial.Open(options); err != nil {
		err = errors.Wrapf(err, "Could not open serial port")
		return
	}

	logger.Infof(" make new channels \n")
	// long lived reader goroutine so we retain state of the stream across individual "reads"
	readerChannel = make(chan serialReadResponse)
	logger.Infof(" make new goroutine \n")
	go func(s SerialIo) {
		defer s.stream.Close()

		var arr []byte = make([]byte, 1)
		var emptyCount = 0
		var sleepCount = 0
		var EMPTY_PER_SLEEP = 5
		for {
			var response serialReadResponse
			var n int
			n, response.err = s.stream.Read(arr)
			if response.err != nil {
				sleepCount = 0
				emptyCount = 0
				readerChannel <- response
			} else if n == 1 {
				if (emptyCount + sleepCount*EMPTY_PER_SLEEP) > 0 {
					logger.Infof("got %d empties before this read\n", emptyCount+sleepCount*EMPTY_PER_SLEEP)
				}
				sleepCount = 0
				emptyCount = 0
				response.data = arr[0]
				readerChannel <- response
			} else {
				emptyCount = emptyCount + 1

				if emptyCount > EMPTY_PER_SLEEP {
					// HACK: on windows, despite asking for blocking IO
					// the Read is returning immediately with
					// n == 0, but no error.  Sleep for a
					// while so we don't chew up infinite CPU.
					// However, if we sleep each time, we get REALLY SLOW IO.
					//
					// FIXME: I dont like picking magic numbers to tune performance
					// need to fix the underlying serial library instead.
					sleepCount = sleepCount + 1
					time.Sleep(time.Duration(10) * time.Millisecond)
					emptyCount = 0
				}
			}
		}
	}(s)

	return
}

func (s SerialIo) close() (err error) {
	logger.Infof(" --> serial.close()\n")
	logger.Infof(" close....1\n")
	if err = s.stream.Close(); err != nil {
		logger.Infof(" close....2\n")
		return
	}
	logger.Infof(" close....done\n")
	return
}

func (s SerialIo) readByte(timeoutMS uint) (b byte, err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	select {
	case response := <-readerChannel:
		if response.err != nil {
			return response.data, errors.Wrap(err, "failed to read byte")
		}
		return response.data, nil
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		return 0, errors.Errorf("TIMEOUT: read timed out at %d ms", timeoutMS)
	}
}

func (s SerialIo) readBytes(timeoutMS uint, num_bytes uint16) (bytes []byte, err error) {
	var arr []byte = make([]byte, num_bytes)

	for i := uint16(0); i < num_bytes; i++ {
		if arr[i], err = s.readByte(timeoutMS); err != nil {
			bytes = arr[0:i]
			err = errors.Wrap(err, "failed to read all bytes")
			return
		}
	}
	bytes = arr
	return
}

func (s SerialIo) writeByte(timeoutMS uint, b byte) (err error) {
	var arr []byte = make([]byte, 1)
	arr[0] = b

	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		_, writeerr := s.stream.Write(arr)
		c <- writeerr
	}()

	select {
	case err := <-c:
		if err != nil {
			return errors.Wrap(err, "failed to write byte")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}

func (s SerialIo) writeBytes(timeoutMS uint, arr []byte) (err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	c := make(chan error, 1)
	go func() {
		_, writeerr := s.stream.Write(arr)
		c <- writeerr
	}()

	select {
	case err := <-c:
		if err != nil {
			return errors.Wrap(err, "failed to write bytes")
		}
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}
