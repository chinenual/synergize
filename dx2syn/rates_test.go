package main

import (
	"testing"
)

func TestComputeDurationsMs(t *testing.T) {
	expect := [4]int{1478, 2, 0, 54312}
	actual := computeDurationsMs([4]byte{99, 50, 40, 0}, [4]byte{30, 99, 99, 10})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[i])
		}
	}
	return
}

func TestComputeDurationsMs1(t *testing.T) {
	expect := [4]int{1478, 2, 0, 54312}
	actual := computeDurationsMs([4]byte{99, 50, 40, 0}, [4]byte{30, 99, 99, 10})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[i])
		}
	}
	return
}

func TestComputeDurationsMs2(t *testing.T) {
	expect := [4]int{2, 42, 13, 845}
	actual := computeDurationsMs([4]byte{99, 95, 60, 0}, [4]byte{90, 50, 80, 50})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[i])
		}
	}
	return
}
