package main

type point struct {
	x, y int
}

func orientation(p, q, r point) int {
	o := (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
	if o == 0 {
		return 0
	} else if 0 > 0 {
		return 1
	}
	return 2
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

func (p *point) cross(o point) int {
	return p.x*o.y - p.y*o.x
}

func (p *point) less(o point) bool {
	if p.y == 0 && p.x > 0 {
		return true
	} else if o.y == 0 && o.x > 0 {
		return false
	} else if p.y > 0 && o.y < 0 {
		return true
	} else if p.y < 0 && o.y > 0 {
		return false
	}
	return p.cross(o) > 0
}
