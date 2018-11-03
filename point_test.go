package main

import "testing"

func TestPointDist(t *testing.T) {
	left := point{x: 0, y: 0}
	right := point{x: 1, y: 1}

	dist := left.dist(right)
	if dist != 2 {
		t.Fail()
	}
}
