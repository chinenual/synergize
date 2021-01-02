package main

import (
	"testing"
)

func TestComputeDurationsMs(t *testing.T) {
	expect := [4]int{6, 7, 7, 45}
	actual := computeDurationsMs([4]int{99, 50, 40, 0}, [4]int{30, 99, 99, 10})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[0])
		}
	}
	return
}

func TestComputeDurationsMs1(t *testing.T) {
	expect := [4]int{1478, 2, 0, 54312}
	actual := computeDurationsMs([4]int{99, 50, 40, 0}, [4]int{30, 99, 99, 10})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[0])
		}
	}
	return
}

func TestComputeDurationsMs2(t *testing.T) {
	expect := [4]int{2, 41, 13, 844}
	actual := computeDurationsMs([4]int{99, 95, 60, 0}, [4]int{90, 50, 80, 50})
	for i, exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[0])
		}
	}
	return
}
