package main

import "github.com/nsf/termbox-go"

type sprite struct {
	x, y int
	c    rune
}

func (s *sprite) move(dx, dy int, l *level) {
	width, height := termbox.Size()
	x := cap(0, width, s.x+dx)
	y := cap(0, height, s.y+dy)

	if l.tile(x, y).isPassable(s) {
		s.x = x
		s.y = y
	}
}

func (s *sprite) handleKeyEvent(e termbox.Event, l *level) {
	if e.Ch != 0 {
		return
	}
	switch e.Key {
	case termbox.KeyArrowLeft:
		s.move(-1, 0, l)
		return
	case termbox.KeyArrowRight:
		s.move(1, 0, l)
		return
	case termbox.KeyArrowUp:
		s.move(0, -1, l)
		return
	case termbox.KeyArrowDown:
		s.move(0, 1, l)
		return
	}
}

func (s *sprite) draw() {
	drawString(s.x, s.y, string(s.c))
}
