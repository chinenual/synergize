package data

import (
	"testing"
)

func AssertByte(t *testing.T, b byte, expected byte, context string) {
	if b != expected {
		t.Errorf("expected %s %04x, got %04x\n", context, expected, b)
	}
}

func AssertUint16(t *testing.T, b uint16, expected uint16, context string) {
	if b != expected {
		t.Errorf("expected %s %04x, got %04x\n", context, expected, b)
	}
}

