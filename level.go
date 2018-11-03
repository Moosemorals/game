package main

import "github.com/nsf/termbox-go"

type attr uint32

const (
	visited attr = 1 << iota
)

func (a *attr) set(n attr) {
	*a = *a | n
}

func (a *attr) has(n attr) bool {
	return *a&n != 0
}

type level struct {
	layout        []tiler
	attributes    []attr
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
				if x == left {
					if y == top {
						l.setTile(x, y, &wall{g: wallTopLeft})
					} else if y == bottom {
						l.setTile(x, y, &wall{g: wallBottomLeft})
					} else {
						l.setTile(x, y, &wall{g: wallVertical})
					}
				} else if x == right {
					if y == top {
						l.setTile(x, y, &wall{g: wallTopRight})
					} else if y == bottom {
						l.setTile(x, y, &wall{g: wallBottomRight})
					} else {
						l.setTile(x, y, &wall{g: wallVertical})
					}
				} else {
					l.setTile(x, y, &wall{g: wallHorizontal})
				}
			} else {
				l.setTile(x, y, new(floor))
			}
		}
	}
	l.setTile(right, (bottom-top)/2, &door{
		open:       false,
		horizontal: false,
		point: point{
			x: right,
			y: (bottom - top) / 2,
		},
	})
}

func (l *level) draw() {
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			var c glyph
			tile := l.tile(x, y)
			if tile != nil && l.hasVisited(x, y) {
				c = tile.glyph()
			} else {
				c = ' '
			}
			drawString(x, y, string(c))
		}
	}
}

func makeLevel(w, h int) *level {
	var l = level{
		width:      w,
		height:     h,
		layout:     make([]tiler, w*h),
		attributes: make([]attr, w*h),
	}

	l.drawRoom(1, 1, 8, 8)
	return &l
}

func (l *level) handleKeyEvent(e termbox.Event, context *context) {
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			tile := l.tile(x, y)
			handler, ok := tile.(keyHandler)
			if ok {
				handler.handleKeyEvent(e, context)
			}
		}
	}
}

func (l *level) isValidPoint(x, y int) bool {
	index := y*l.width + x
	return index >= 0 && index < len(l.layout)
}

func (l *level) setAttribute(x, y int, a attr) {
	l.attributes[y*l.width+x].set(visited)
}

func (l *level) visit(x, y int) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if l.isValidPoint(x+dx, y+dy) {
				l.setAttribute(x+dx, y+dy, visited)
			}
		}
	}
}

func (l *level) hasVisited(x, y int) bool {
	return l.attributes[y*l.width+x].has(visited)
}
