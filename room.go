package main

import (
	"math"
	"math/rand"
)

func random(min, max int) int {
	if (max - min) <= 0 {
		return 0
	}
	return rand.Intn(max-min) + min
}

type room struct {
	x, y, w, h int
}

type node struct {
	box         *room
	left, right *node
}

func (n *node) rooms() []*room {
	if n.left != nil && n.right != nil {
		return append(n.left.rooms(), n.right.rooms()...)
	}
	return []*room{n.box}
}

func (r *room) center() point {
	return point{r.x + r.w/2, r.y + r.h/2}
}

func (r *room) link(d *room, l *level) {
	var start, end point
	if d.x+d.w+1 > r.x && d.x < r.x+r.w-1 {
		// Horizontal overlap
		if r.x >= d.x {
			start.x = r.x + (((d.x + d.w) - r.x) / 2)
		} else {
			start.x = d.x + (((r.x + r.w) - d.x) / 2)
		}
		start.y = r.y + r.h - 1
		end.x = start.x
		end.y = d.y + 1
		l.drawCorridor(start, end)
	}
	if d.y+d.h+1 > r.y && d.y < r.y+r.h-1 {
		// vertical overlap
		if r.y >= d.y {
			start.y = r.y + (((d.y + d.h) - r.y) / 2)
		} else {
			start.y = d.y + (((r.y + r.h) - d.y) / 2)
		}
		start.x = r.x + r.w - 1
		end.y = start.y
		end.x = d.x + 1
		l.drawCorridor(start, end)
	}
}

const (
	overlapHorizontal int = 1 << iota
	overlapVertical
)

func (r *room) overlaps(d *room) int {
	return d.x
}

func splitRoom(r *room, vh bool) (left, right *room) {
	var split int
	if vh {
		factor := float64(r.w) / (rand.Float64() + 1.5)
		split = int(math.Round(factor))

		// Vertical split
		left = &room{
			x: r.x,
			y: r.y,
			w: split,
			h: r.h,
		}
		right = &room{
			x: r.x + left.w,
			y: r.y,
			w: r.w - left.w,
			h: r.h,
		}
	} else {
		factor := float64(r.h) / (rand.Float64() + 1.5)
		split = int(math.Round(factor))

		// Horizontal split
		left = &room{
			x: r.x,
			y: r.y,
			w: r.w,
			h: split,
		}
		right = &room{
			x: r.x,
			y: r.y + left.h,
			w: r.w,
			h: r.h - left.h,
		}
	}
	return
}

func buildRooms(r *room, depth int) *node {
	root := &node{box: r}
	if depth > 0 {
		leftBox, rightBox := splitRoom(r, depth%2 == 0)
		root.left = buildRooms(leftBox, depth-1)
		root.right = buildRooms(rightBox, depth-1)
	}
	return root
}

func resizeRoom(r *room) *room {
	var out room
	tries := 3
	for {
		out.x = r.x + random(1, r.x/3)
		out.y = r.y + random(1, r.y/3)
		out.w = r.w - (out.x - r.x)
		out.h = r.h - (out.y - r.y)
		if out.w > 3 && out.h > 3 {
			return &out
		}
		tries--
		if tries == 0 {
			return nil
		}
	}
}

func (r *room) draw(l *level) {
	var p point
	top := r.y
	left := r.x
	right := r.x + r.w - 1
	bottom := r.y + r.h - 1
	for p.x = left; p.x <= right; p.x++ {
		for p.y = top; p.y <= bottom; p.y++ {
			if p.x == left || p.x == right || p.y == top || p.y == bottom {
				if p.x == left {
					if p.y == top {
						l.setTile(p, &wall{g: wallTopLeft})
					} else if p.y == bottom {
						l.setTile(p, &wall{g: wallBottomLeft})
					} else {
						l.setTile(p, &wall{g: wallVertical})
					}
				} else if p.x == right {
					if p.y == top {
						l.setTile(p, &wall{g: wallTopRight})
					} else if p.y == bottom {
						l.setTile(p, &wall{g: wallBottomRight})
					} else {
						l.setTile(p, &wall{g: wallVertical})
					}
				} else {
					l.setTile(p, &wall{g: wallHorizontal})
				}
			} else {
				l.setTile(p, new(floor))
			}
		}
	}
}
