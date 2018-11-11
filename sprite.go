package main

import (
	"github.com/nsf/termbox-go"
)

type sprite struct {
	point
	c glyph
}

func (s *sprite) move(dx, dy int, c *context) {
	size := c.size
	dest := point{
		x: cap(0, size.x, s.x+dx),
		y: cap(0, size.y, s.y+dy),
	}
	tile := c.level.tile(dest)
	if tile != nil && tile.isPassable(s) {
		s.point = dest
		c.level.visit(dest)
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
	drawGlyph(s.point, s.c, termbox.ColorWhite, termbox.ColorBlack)
}
