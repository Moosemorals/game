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

func (l *level) toIndex(p point) int {
	return p.y*l.width + p.x
}

func (l *level) setTile(p point, tile tiler) {
	l.layout[l.toIndex(p)] = tile
}

func (l *level) tile(p point) tiler {
	return l.layout[l.toIndex(p)]
}

func (l *level) drawRoom(top, left, bottom, right int) {
	var p point
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
	var dpoint = point{right, (bottom - top) / 2}
	l.setTile(dpoint, &door{
		open:       false,
		horizontal: false,
		point:      dpoint,
	})
}

func (l *level) draw() {
	var p point
	for p.x = 0; p.x < l.width; p.x++ {
		for p.y = 0; p.y < l.height; p.y++ {
			var c glyph
			tile := l.tile(p)
			if tile != nil && l.hasVisited(p) {
				c = tile.glyph()
			} else {
				c = ' '
			}
			drawString(p, string(c))
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
	var p point
	for p.x = 0; p.x < l.width; p.x++ {
		for p.y = 0; p.y < l.height; p.y++ {
			tile := l.tile(p)
			handler, ok := tile.(keyHandler)
			if ok {
				handler.handleKeyEvent(e, context)
			}
		}
	}
}

func (l *level) isValidPoint(p point) bool {
	index := p.y*l.width + p.x
	return index >= 0 && index < len(l.layout)
}

func (l *level) setAttribute(p point, a attr) {
	l.attributes[p.y*l.width+p.x].set(visited)
}

func (l *level) visit(p point) {
	var delta point
	for delta.x = -1; delta.x <= 1; delta.x++ {
		for delta.y = -1; delta.y <= 1; delta.y++ {
			if l.isValidPoint(p.add(delta)) {
				l.setAttribute(p.add(delta), visited)
			}
		}
	}
}

func (l *level) hasVisited(p point) bool {
	return l.attributes[p.y*l.width+p.x].has(visited)
}
