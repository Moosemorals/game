package main

import "github.com/nsf/termbox-go"

type level struct {
	layout        []tiler
	width, height int
}

func (l *level) setTile(x, y int, tile tiler) {
	l.layout[y*l.width+x] = tile
}

func (l *level) tile(x, y int) tiler {
	return l.layout[y*l.width+x]
}

func (l *level) drawRoom(top, left, bottom, right int) {
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			if x == left || x == right || y == top || y == bottom {
				l.setTile(x, y, new(wall))
			} else {
				l.setTile(x, y, new(floor))
			}
		}
	}
}

func (l *level) draw() {
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			var c tile
			tile := l.tile(x, y)
			if tile == nil {
				c = ' '
			} else {
				c = tile.tile()
			}
			drawString(x, y, string(c))
		}
	}
}

func makeLevel() level {
	w, h := termbox.Size()
	var l = level{
		width:  w,
		height: h,
		layout: make([]tiler, w*h),
	}

	l.drawRoom(1, 1, 8, 8)
	return l
}