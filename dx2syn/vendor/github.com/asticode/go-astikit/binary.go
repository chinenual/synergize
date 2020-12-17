package astikit

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// BitsWriter represents an object that can write individual bits into a writer
// in a developer-friendly way. Check out the Write method for more information.
// This is particularly helpful when you want to build a slice of bytes based
// on individual bits for testing purposes.
type BitsWriter struct {
	bo binary.ByteOrder
	rs []rune
	w  io.Writer
}

// BitsWriterOptions represents BitsWriter options
type BitsWriterOptions struct {
	ByteOrder binary.ByteOrder
	Writer    io.Writer
}

// NewBitsWriter creates a new BitsWriter
func NewBitsWriter(o BitsWriterOptions) (w *BitsWriter) {
	w = &BitsWriter{
		bo: o.ByteOrder,
		w:  o.Writer,
	}
	if w.bo == nil {
		w.bo = binary.BigEndian
	}
	return
}

// Write writes bits into the writer. Bits are only written when there are
// enough to create a byte. When using a string or a bool, bits are added
// from left to right as if
// Available types are:
//   - string("10010"): processed as n bits, n being the length of the input
//   - []byte: processed as n bytes, n being the length of the input
//   - bool: processed as one bit
//   - uint8/uint16/uint32/uint64: processed as n bits, if type is uintn
func (w *BitsWriter) Write(i interface{}) (err error) {
	// Transform input into "10010" format
	var s string
	switch a := i.(type) {
	case string:
		s = a
	case []byte:
		for _, b := range a {
			s += fmt.Sprintf("%.8b", b)
		}
	case bool:
		if a {
			s = "1"
		} else {
			s = "0"
		}
	case uint8:
		s = fmt.Sprintf("%.8b", i)
	case uint16:
		s = fmt.Sprintf("%.16b", i)
	case uint32:
		s = fmt.Sprintf("%.32b", i)
	case uint64:
		s = fmt.Sprintf("%.64b", i)
	default:
		err = errors.New("astikit: invalid type")
		return
	}

	// Loop through runes
	for _, r := range s {
		// Append rune
		if w.bo == binary.LittleEndian {
			w.rs = append([]rune{r}, w.rs...)
		} else {
			w.rs = append(w.rs, r)
		}

		// There are enough bits to form a byte
		if len(w.rs) == 8 {
			// Get value
			v := w.val()

			// Remove runes
			w.rs = w.rs[8:]

			// Write
			if err = binary.Write(w.w, w.bo, v); err != nil {
				return
			}
		}
	}
	return
}

func (w *BitsWriter) val() (v uint8) {
	var power float64
	for idx := len(w.rs) - 1; idx >= 0; idx-- {
		if w.rs[idx] == '1' {
			v = v + uint8(math.Pow(2, power))
		}
		power++
	}
	return
}

// WriteN writes the input into n bits
func (w *BitsWriter) WriteN(i interface{}, n int) error {
	switch i.(type) {
	case uint8, uint16, uint32, uint64:
		return w.Write(fmt.Sprintf(fmt.Sprintf("%%.%db", n), i))
	default:
		return errors.New("astikit: invalid type")
	}
}

var byteHamming84Tab = [256]uint8{
	0x01, 0xff, 0xff, 0x08, 0xff, 0x0c, 0x04, 0xff, 0xff, 0x08, 0x08, 0x08, 0x06, 0xff, 0xff, 0x08,
	0xff, 0x0a, 0x02, 0xff, 0x06, 0xff, 0xff, 0x0f, 0x06, 0xff, 0xff, 0x08, 0x06, 0x06, 0x06, 0xff,
	0xff, 0x0a, 0x04, 0xff, 0x04, 0xff, 0x04, 0x04, 0x00, 0xff, 0xff, 0x08, 0xff, 0x0d, 0x04, 0xff,
	0x0a, 0x0a, 0xff, 0x0a, 0xff, 0x0a, 0x04, 0xff, 0xff, 0x0a, 0x03, 0xff, 0x06, 0xff, 0xff, 0x0e,
	0x01, 0x01, 0x01, 0xff, 0x01, 0xff, 0xff, 0x0f, 0x01, 0xff, 0xff, 0x08, 0xff, 0x0d, 0x05, 0xff,
	0x01, 0xff, 0xff, 0x0f, 0xff, 0x0f, 0x0f, 0x0f, 0xff, 0x0b, 0x03, 0xff, 0x06, 0xff, 0xff, 0x0f,
	0x01, 0xff, 0xff, 0x09, 0xff, 0x0d, 0x04, 0xff, 0xff, 0x0d, 0x03, 0xff, 0x0d, 0x0d, 0xff, 0x0d,
	0xff, 0x0a, 0x03, 0xff, 0x07, 0xff, 0xff, 0x0f, 0x03, 0xff, 0x03, 0x03, 0xff, 0x0d, 0x03, 0xff,
	0xff, 0x0c, 0x02, 0xff, 0x0c, 0x0c, 0xff, 0x0c, 0x00, 0xff, 0xff, 0x08, 0xff, 0x0c, 0x05, 0xff,
	0x02, 0xff, 0x02, 0x02, 0xff, 0x0c, 0x02, 0xff, 0xff, 0x0b, 0x02, 0xff, 0x06, 0xff, 0xff, 0x0e,
	0x00, 0xff, 0xff, 0x09, 0xff, 0x0c, 0x04, 0xff, 0x00, 0x00, 0x00, 0xff, 0x00, 0xff, 0xff, 0x0e,
	0xff, 0x0a, 0x02, 0xff, 0x07, 0xff, 0xff, 0x0e, 0x00, 0xff, 0xff, 0x0e, 0xff, 0x0e, 0x0e, 0x0e,
	0x01, 0xff, 0xff, 0x09, 0xff, 0x0c, 0x05, 0xff, 0xff, 0x0b, 0x05, 0xff, 0x05, 0xff, 0x05, 0x05,
	0xff, 0x0b, 0x02, 0xff, 0x07, 0xff, 0xff, 0x0f, 0x0b, 0x0b, 0xff, 0x0b, 0xff, 0x0b, 0x05, 0xff,
	0xff, 0x09, 0x09, 0x09, 0x07, 0xff, 0xff, 0x09, 0x00, 0xff, 0xff, 0x09, 0xff, 0x0d, 0x05, 0xff,
	0x07, 0xff, 0xff, 0x09, 0x07, 0x07, 0x07, 0xff, 0xff, 0x0b, 0x03, 0xff, 0x07, 0xff, 0xff, 0x0e,
}

// ByteHamming84Decode hamming 8/4 decodes
func ByteHamming84Decode(i uint8) (o uint8, ok bool) {
	o = byteHamming84Tab[i]
	if o == 0xff {
		return
	}
	ok = true
	return
}

var byteParityTab = [256]uint8{
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00,
}

// ByteParity returns the byte parity
func ByteParity(i uint8) (o uint8, ok bool) {
	ok = byteParityTab[i] == 1
	o = i & 0x7f
	return
}
