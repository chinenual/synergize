package io

import (
	"log"
)

type IoImpl interface {
	readByte(timeoutMS uint) (b byte, err error)
	readBytes(timeoutMS uint, num_bytes uint16) (bytes []byte, err error)
	writeByte(timeoutMS uint, b byte) (err error)
	writeBytes(timeoutMS uint, arr []byte) (err error)
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

// HACK: using ugly WithTimeout suffix to keep govet from complaining that ReadByte (et.al.)
// have different signatures than "standard" Read/Write byte methods
func (c *Conn) ReadByteWithTimeout(timeoutMS uint, purpose string) (b byte, err error) {

	if c.verbose {
		log.Printf("       serial.Read (%d ms) - %s\n", timeoutMS, purpose)
	}

	b, err = c.impl.readByte(timeoutMS)
	if c.verbose {
		if err != nil {
			log.Printf("       read err: %v\n", err)
		} else {
			log.Printf(" %02x <-- serial.Read (%v ms)\n", b, timeoutMS)
		}
	}
	if c.record {
		c.recordIn = append(c.recordIn, b)
	}
	return
}

func (c *Conn) ReadBytesWithTimeout(timeoutMS uint, num_bytes uint16, purpose string) (bytes []byte, err error) {
	if c.verbose {
		log.Printf("       serial.Read %d bytes (%d ms) - %s\n", num_bytes, timeoutMS, purpose)
	}
	bytes, err = c.impl.readBytes(timeoutMS, num_bytes)
	if c.verbose {
		if err != nil {
			log.Printf("       read err: %v\n", err)
		} else {
			log.Printf(" %02x <-- serial.Read (%d ms)\n", bytes, timeoutMS)
		}
	}
	if c.record {
		c.recordIn = append(c.recordIn, bytes...)
	}
	return
}

func (c *Conn) WriteByteWithTimeout(timeoutMS uint, b byte, purpose string) (err error) {
	if c.verbose {
		log.Printf(" --> %02x serial.Write (%d ms) - %s\n", b, timeoutMS, purpose)
	}

	err = c.impl.writeByte(timeoutMS, b)
	if c.verbose && err != nil {
		log.Printf("        write err: %v\n", err)
	}
	if c.record {
		c.recordOut = append(c.recordOut, b)
	}
	return
}

func (c *Conn) WriteBytesWithTimeout(timeoutMS uint, arr []byte, purpose string) (err error) {
	if c.verbose {
		log.Printf(" --> %02x serial.WriteBytes (%d ms) - %s\n", arr, timeoutMS, purpose)
	}
	err = c.impl.writeBytes(timeoutMS, arr)
	if c.verbose && err != nil {
		log.Printf("        write err: %v\n", err)
	}
	if c.record {
		c.recordOut = append(c.recordOut, arr...)
	}
	return
}