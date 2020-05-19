package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	var err error
	var bytes []byte
	if bytes, err = ioutil.ReadFile("./VERSION"); err != nil {
		t.Fatalf("Error when reading VERSION: %v", err)
	}
	fromFile := strings.TrimSpace(strings.Split(string(bytes), "=")[1])

	if Version != fromFile {
		t.Fatalf("Expected '%s', got '%s'", Version, fromFile)
	}

}
