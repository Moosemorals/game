package main

import (
	"github.com/nsf/termbox-go"
)

type sprite struct {
	point
	c rune
}

func (s *sprite) move(dx, dy int, context *context) {
	width := context.size.x
	height := context.size.y

	x := cap(0, width, s.x+dx)
	y := cap(0, height, s.y+dy)

	tile := context.level.tile(x, y)
	if tile != nil && tile.isPassable(s) {
		s.x = x
		s.y = y
	}
}

func (s *sprite) handleKeyEvent(e termbox.Event, context *context) {
	if e.Ch != 0 {
		return
	}
	switch e.Key {
	case termbox.KeyArrowLeft:
		s.move(-1, 0, context)
		return
	case termbox.KeyArrowRight:
		s.move(1, 0, context)
		return
	case termbox.KeyArrowUp:
		s.move(0, -1, context)
		return
	case termbox.KeyArrowDown:
		s.move(0, 1, context)
		return
	}
}

func (s *sprite) draw() {
	drawString(s.x, s.y, string(s.c))
}
