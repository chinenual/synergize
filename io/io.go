package io

import (
	"log"
)

type IoImpl interface {
	// HACK: using ugly WithTimeout suffix to keep govet from complaining that ReadByte (et.al.)
	// have different signatures than "standard" Read/Write byte methods
	ReadByteWithTimeout(timeoutMS uint) (b byte, err error)
	ReadBytesWithTimeout(timeoutMS uint, num_bytes uint16) (bytes []byte, err error)
	WriteByteWithTimeout(timeoutMS uint, b byte) (err error)
	WriteBytesWithTimeout(timeoutMS uint, arr []byte) (err error)
}

type Conn struct {
	verbose   bool
	record    bool
	recordIn  []byte
	recordOut []byte
	impl      IoImpl
}

func IoInit(impl IoImpl, verbose bool) (c Conn, err error) {
	c.verbose = verbose
	c.impl = impl
	return
}

func (c *Conn) StartRecord() {
	c.record = true
	c.recordIn = []byte{}
	c.recordOut = []byte{}
}

func (c *Conn) CloseRecord() {
	c.record = false
	c.recordIn = nil
	c.recordOut = nil
}

func (c *Conn) GetRecord() (in, out []byte) {
	in = c.recordIn
	out = c.recordOut
	return
}

func (c *Conn) LoggedReadByteWithTimeout(timeoutMS uint, purpose string) (b byte, err error) {

	if c.verbose {
		log.Printf("       serial.Read (%d ms) - %s\n", timeoutMS, purpose)
	}

	b, err = c.impl.ReadByteWithTimeout(timeoutMS)
	if c.verbose {
		if err != nil {
			log.Printf("       read err: %v\n", err)
		} else {
			log.Printf(" %02x <-- serial.Read (%v ms)\n", b, timeoutMS)
		}
	}
	if c.record {
		log.Printf("before append %d\n", len(c.recordIn))
		c.recordIn = append(c.recordIn, b)
		log.Printf("after append %d\n", len(c.recordIn))
	}
	return
}

func (c *Conn) LoggedReadBytesWithTimeout(timeoutMS uint, num_bytes uint16, purpose string) (bytes []byte, err error) {
	if c.verbose {
		log.Printf("       serial.Read %d bytes (%d ms) - %s\n", num_bytes, timeoutMS, purpose)
	}
	bytes, err = c.impl.ReadBytesWithTimeout(timeoutMS, num_bytes)
	if c.verbose {
		if err != nil {
			log.Printf("       read err: %v\n", err)
		} else {
			log.Printf(" %02x <-- serial.Read (%d ms)\n", bytes, timeoutMS)
		}
	}
	if c.record {
		log.Printf("before append %d\n", len(c.recordIn))
		c.recordIn = append(c.recordIn, bytes...)
		log.Printf("after append %d\n", len(c.recordIn))
	}
	return
}

func (c *Conn) LoggedWriteByteWithTimeout(timeoutMS uint, b byte, purpose string) (err error) {
	if c.verbose {
		log.Printf(" --> %02x serial.Write (%d ms) - %s\n", b, timeoutMS, purpose)
	}

	err = c.impl.WriteByteWithTimeout(timeoutMS, b)
	if c.verbose && err != nil {
		log.Printf("        write err: %v\n", err)
	}
	if c.record {
		log.Printf("before append %d\n", len(c.recordOut))
		c.recordOut = append(c.recordOut, b)
		log.Printf("after append %d\n", len(c.recordOut))
	}
	return
}

func (c *Conn) LoggedWriteBytesWithTimeout(timeoutMS uint, arr []byte, purpose string) (err error) {
	if c.verbose {
		log.Printf(" --> %02x serial.WriteBytes (%d ms) - %s\n", arr, timeoutMS, purpose)
	}
	err = c.impl.WriteBytesWithTimeout(timeoutMS, arr)
	if c.verbose && err != nil {
		log.Printf("        write err: %v\n", err)
	}
	if c.record {
		log.Printf("before append %d\n", len(c.recordOut))
		c.recordOut = append(c.recordOut, arr...)
		log.Printf("after append %d\n", len(c.recordOut))
	}
	return
}
