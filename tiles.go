package main

import (
	"github.com/nsf/termbox-go"
)

type glyph rune

type tiler interface {
	isPassable(s *sprite) bool
	glyph() glyph
}

const (
	wallTopLeft     glyph = '\u250f'
	wallHorizontal  glyph = '\u2501'
	wallTopRight    glyph = '\u2513'
	wallVertical    glyph = '\u2503'
	wallBottomLeft  glyph = '\u2517'
	wallBottomRight glyph = '\u251b'
	floorTile       glyph = '.'
)

type wall struct {
	g glyph
}

func (w *wall) isPassable(s *sprite) bool {
	return false
}

func (w *wall) glyph() glyph {
	return w.g
}

type floor struct{}

func (f *floor) isPassable(s *sprite) bool {
	return true
}

func (f *floor) glyph() glyph {
	return floorTile
}

type door struct {
	point
	open       bool
	horizontal bool
}

func (d *door) isPassable(s *sprite) bool {
	return d.open
}

func (d *door) glyph() glyph {
	if d.open {
		return floorTile
	} else if d.horizontal {
		return '\u2500'
	}
	return '\u2502'
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
