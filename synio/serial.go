package synio

import (
	"fmt"
	"io"
	"time"
	"log"
	"runtime"
	"github.com/pkg/errors"
	"github.com/jacobsa/go-serial/serial"
)

var serialVerbose = false
var readerChannel chan serialReadResponse
var readerChannelQuit chan bool

type serialReadResponse struct {
	data byte
	err error
}

var (
       stream io.ReadWriteCloser
)

func serialInit(port string, baudRate uint, verbose bool) (err error) {
	serialVerbose = verbose
	options := serial.OpenOptions{
		PortName: port,
		BaudRate: baudRate,
		ParityMode: serial.PARITY_NONE,
		RTSCTSFlowControl: true,
		InterCharacterTimeout: 500,
		MinimumReadSize: 1,
		DataBits: 8,
		StopBits: 1,
	}
	if readerChannel != nil {
		if runtime.GOOS == "darwin" {
			// FIXME: can't find a way to interrupt the blocking
			// read in the goroutine; until that's changed, just
			// throw an error if user tries to reopen the serial port
			//
			// This "works" on Windows due to the bug that is causing
			// the Read to be non-blocking
			err = errors.New("Cannot re-open the serial connection.  To reinitialize with a new port or baud rate you must restart the Synergize application.")
			return
		}
		readerChannelQuit <- true
	}
	log.Printf(" --> serial.Open(%#v)\n",options)
	if stream,err = serial.Open(options); err != nil {
		return errors.Wrapf(err,"Could not open serial port")
	}

	log.Printf(" make new channels \n")
	// long lived reader goroutine so we retain state of the stream across individual "reads"
	readerChannel = make(chan serialReadResponse)
	readerChannelQuit = make(chan bool)
	log.Printf(" make new goroutine \n")
	go func () {
		defer stream.Close()
		
		var arr []byte = make([]byte,1);
		var emptyCount = 0
		var sleepCount = 0
		var EMPTY_PER_SLEEP = 5
		for {
			select {
			case <- readerChannelQuit:
				log.Printf(" closing serial channel\n")
				close(readerChannelQuit)
				close(readerChannel)
				log.Printf(" ending goroutine\n")
				return				
			default:
				var response serialReadResponse
				var n int
				n, response.err = stream.Read(arr);
				if response.err != nil {
					sleepCount = 0
					emptyCount = 0
					// if response.data==0  && serialVerbose {
					//	 log.Printf("   **** error %v\n",err)
					// }
					readerChannel <- response
				} else if n == 1 {
					if (emptyCount + sleepCount*EMPTY_PER_SLEEP)>0 {
						log.Printf("got %d empties before this read\n",emptyCount + sleepCount*EMPTY_PER_SLEEP)
					}
					sleepCount = 0
					emptyCount = 0
					response.data = arr[0]
					// if response.data==0  && serialVerbose {
					//	 log.Printf("   **** Zero w/ n=1 and no error\n")
					// }
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
		}
	}()
	
	return nil
}
	
func serialReadByte(timeoutMS uint, purpose string) (b byte, err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if serialVerbose {log.Printf("       serial.Read (%d ms) - %s\n",timeoutMS,purpose)}
	
	select {
	case response := <-readerChannel:
		if response.err != nil {
			if serialVerbose {log.Printf(" %v <-- serial.Read (%d ms)\n",response.err,timeoutMS)}
			return response.data,errors.Wrap(err, "failed to read byte")
		}
		if serialVerbose {log.Printf(" %02x <-- serial.Read (%d ms)\n",response.data,timeoutMS)}
		return response.data, nil
	case <-time.After(time.Millisecond * time.Duration(timeoutMS)):
		// call timed out
		if serialVerbose {log.Printf("   read TIMEOUT at %d ms (%x)\n", timeoutMS,0)}
		return 0,errors.Errorf("TIMEOUT: read timed out at %d ms", timeoutMS)
	}
	if serialVerbose {log.Printf("   Unexpected read error\n")}
	return 0, errors.Errorf("Unexpected read error")
}


func serialReadBytes(timeoutMS uint, num_bytes uint16, purpose string) (bytes []byte, err error) {
	var arr []byte = make([]byte,num_bytes);

	for i:= uint16(0); i < num_bytes; i++ {
		if arr[i],err = serialReadByte(timeoutMS, fmt.Sprintf("%s: %d",purpose,i)); err != nil {
			bytes = arr[0:i]
			err = errors.Wrap(err, "failed to read all bytes")
			return
		}
	}
	bytes = arr
	return
}
	

func serialWriteByte(timeoutMS uint, b byte, purpose string) (err error) {
	var arr []byte = make([]byte,1);
	arr[0] = b;
	
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if serialVerbose {log.Printf(" --> %02x serial.Write (%d ms) - %s\n",arr[0], timeoutMS,purpose)}

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
		if serialVerbose {log.Printf("   write TIMEOUT at %d ms\n", timeoutMS)}
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}

func serialWriteBytes(timeoutMS uint, arr []byte, purpose string) (err error) {
	// use goroutines to handle timeout of synchronous IO.
	// See https://github.com/golang/go/wiki/Timeouts

	if serialVerbose {log.Printf(" --> %02x serial.WriteBytes (%d ms) - %s\n",arr,timeoutMS,purpose)}

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
		if serialVerbose {log.Printf("   write TIMEOUT at %d ms\n", timeoutMS)}
		return errors.Errorf("TIMEOUT: write timed out at %d ms", timeoutMS)
	}
	return nil
}
