package main

import (
	"github.com/nsf/termbox-go"
)

type tile rune

type tiler interface {
	isPassable(s *sprite) bool
	tile() tile
}

type wall struct {
}

func (w *wall) isPassable(s *sprite) bool {
	return false
}

func (w *wall) tile() tile {
	return '#'
}

type floor struct {
}

func (f *floor) isPassable(s *sprite) bool {
	return true
}

func (f *floor) tile() tile {
	return '.'
}

type door struct {
	point
	open       bool
	horizontal bool
}

func (d *door) isPassable(s *sprite) bool {
	return d.open
}

func (d *door) tile() tile {
	if d.open && d.horizontal {
		return '/'
	} else if d.open && !d.horizontal {
		return '\\'
	} else if !d.open && d.horizontal {
		return '-'
	}
	return '|'
}

func (d *door) handleKeyEvent(e termbox.Event, c *context) {
	if e.Ch == 'o' && d.dist(c.player.point) <= 2 {
		if d.open {
			c.log("Closing door")
		} else {
			c.log("Opening door")
		}
		d.open = !d.open
	}
}
