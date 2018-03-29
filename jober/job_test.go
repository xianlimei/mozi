package jober

import "testing"

func TestJob_Run(t *testing.T) {
	func1 := func() {
		t.Log("func1 run")
	}
	func2 := func() {
		t.Log("func2 run")
	}
	func3 := func() {
		t.Log("func3 run")
	}
	j := &Job{
		ID:    "123",
		execs: []func(){func1, func2, func3},
	}
	j.Run()
}
