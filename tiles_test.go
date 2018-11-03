package main

import "testing"

func TestAttr(t *testing.T) {
	var x attr

	x.set(visited)

	if x != visited {
		t.Fail()
	}

	if !x.has(visited) {
		t.Fail()
	}
}
