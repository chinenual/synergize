package data

import (
	"encoding/json"
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

func (u *ArrayOfByte) UnmarshalJSON(b []byte) (err error) {
	var tmp []int
	// Javascript might write signed bytes to JSON - decode them as int, then convert to byte
	if err = json.Unmarshal(b, &tmp); err != nil {
		return
	}
	var result []byte
	for i := range tmp {
		result = append(result, byte(tmp[i]))
	}
	*u = result
	return nil
}

type SpaceEncodedString [8]byte

func (u *SpaceEncodedString) MarshalJSON() ([]byte, error) {
	result := "\"" + string(u[:]) + "\""
	//fmt.Printf("MARSHAL '%s' -> '%s'\n", u, result)
	return []byte(result), nil
}

func (u *SpaceEncodedString) UnmarshalJSON(s []byte) error {
	// Discard the leading and trailing '""
	s = s[1:(len(s) - 2)]
	for i := range u {
		if i < len(s) {
			u[i] = s[i]
		} else {
			u[i] = ' '
		}
	}
	//fmt.Printf("UNMARSHAL '%s' -> '%s'\n", s, u)
	return nil
}
