package data

import (
	"fmt"
	"strings"
)

type ArrayOfByte []byte

// JSON workarounds for []byte encoding - default go marshaller encodes as Base64
func (u ArrayOfByte) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}
func (u *ArrayOfByte) UnmarshalJSON(b []byte) error {
	if u == nil {
		*u = nil
	} else {
		s := strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
		*u = []byte(s)
	}
	return nil
}

type SpaceEncodedString [8]byte

func (u *SpaceEncodedString) MarshalJSON() ([]byte, error) {
	var result string
	result = "\"" + string(u[:]) + "\""
	fmt.Printf("MARSHAL '%s' -> '%s'\n", u, result)
	return []byte(result), nil
}

func (u *SpaceEncodedString) UnmarshalJSON(s []byte) error {
	// Discard the leading and trailing '""
	s = s[1:(len(s) - 2)]
	for i, _ := range u {
		if i < len(s) {
			u[i] = s[i]
		} else {
			u[i] = ' '
		}
	}
	fmt.Printf("UNMARSHAL '%s' -> '%s'\n", s, u)
	return nil
}
