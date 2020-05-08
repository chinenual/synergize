package data

// JSON workarounds for []byte encoding - default go marshaller encodes as Base64

import (
	"strings"
	"fmt"
)


type ArrayOfByte []byte

type SpaceEncodedString [8]byte

func (u ArrayOfByte) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

func (u SpaceEncodedString) MarshalJSON() ([]byte, error) {
	var result string
	result = "\"" + string(u[:]) + "\""
	return []byte(result), nil
}

