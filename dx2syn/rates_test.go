package main

import (
	"testing"
)

func TestComputeDurationsMs(t *testing.T) {
	expect := [4]int{6, 7, 7, 45}
	actual := computeDurationsMs([4]int{99,80,99,0}, [4]int{80,80,70,80})
	for i,exp := range expect {
		if actual[i] != exp {
			t.Errorf("ms[%d] expected %d, got '%d'\n", i, exp, actual[0])
		}
	}
	return
}
