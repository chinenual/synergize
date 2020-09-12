package io

type MockIo struct {
}

func MockInit() (s MockIo, err error) {
	return
}

func (s MockIo) close() (err error) {
	return
}

func (s MockIo) readByte(timeoutMS uint) (b byte, err error) {
	return
}

func (s MockIo) readBytes(timeoutMS uint, numBytes uint16) (bytes []byte, err error) {
	return
}

func (s MockIo) writeByte(timeoutMS uint, b byte) (err error) {
	return
}

func (s MockIo) writeBytes(timeoutMS uint, arr []byte) (err error) {
	return
}
