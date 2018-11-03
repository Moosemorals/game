package main

type point struct {
	x, y int
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p *point) dist(o point) int {
	return abs(p.x-o.x) + abs(p.y-o.y)
}

func (p *point) add(o point) point {
	return point{p.x + o.x, p.y + o.y}
}
