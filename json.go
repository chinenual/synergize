package main

// JSON workarounds for []uint8 encoding - default go marshaller encodes as Base64

import (
	"strings"
	"fmt"
)


type ArrayOfUint8 []uint8

func (u ArrayOfUint8) MarshalJSON() ([]byte, error) {
    var result string
    if u == nil {
        result = "null"
    } else {
        result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
    }
    return []byte(result), nil
}

