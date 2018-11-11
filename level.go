package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

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
	index := p.y*l.width + p.x

	if index >= len(l.layout) {
		panic(fmt.Sprintf("Well, we're in a pickle\n px %d py %d index %d len %d w %d\n", p.x, p.y, index, len(l.layout), l.width))
	}
	return index
}

func (l *level) setTile(p point, tile tiler) {
	l.layout[l.toIndex(p)] = tile
}

func (l *level) tile(p point) tiler {
	return l.layout[l.toIndex(p)]
}

func (l *level) drawCorridor(start, end point) {
	if start.x == end.x && start.y == end.y {
		log.Panicln("Skipping")
		// Degenerate case, both ends in the same place
		return
	} else if start.x == end.x {
		if start.y < end.y {
			for i := start.y; i < end.y; i++ {
				w := point{start.x, i}
				l.setTile(w, new(floor))
				w.x = start.x - 1
				t := l.tile(w)
				if t == nil {
					l.setTile(w, &wall{wallVertical})
				} else if t.glyph() == wallHorizontal {
					t := l.tile(point{w.x, w.y + 1})
					if t != nil {
						l.setTile(w, &wall{wallBottomRight})
					} else {
						l.setTile(w, &wall{wallTopRight})
					}
				}
				w.x = start.x + 1
				t = l.tile(w)
				if t == nil {
					l.setTile(w, &wall{wallVertical})
				} else if t.glyph() == wallHorizontal {
					t := l.tile(point{w.x, w.y + 1})
					if t != nil {
						l.setTile(w, &wall{wallBottomLeft})
					} else {
						l.setTile(w, &wall{wallTopLeft})
					}
				}
			}
		} else {
			l.drawCorridor(end, start)
		}
	} else if start.y == end.y {
		if start.x < end.x {
			for i := start.x; i < end.x; i++ {
				w := point{i, start.y}
				l.setTile(w, new(floor))
				w.y = start.y - 1
				t := l.tile(w)
				if t == nil {
					l.setTile(w, &wall{wallHorizontal})
				} else if t.glyph() == wallVertical {
					t := l.tile(point{w.x + 1, w.y})
					if t != nil {
						l.setTile(w, &wall{wallBottomRight})
					} else {
						l.setTile(w, &wall{wallBottomLeft})
					}
				}
				w.y = start.y + 1
				t = l.tile(w)
				if t == nil {
					l.setTile(w, &wall{wallHorizontal})
				} else if t.glyph() == wallVertical {
					t := l.tile(point{w.x + 1, w.y})
					if t != nil {
						l.setTile(w, &wall{wallTopRight})
					} else {
						l.setTile(w, &wall{wallTopLeft})
					}

				}
			}
		} else {
			l.drawCorridor(end, start)
		}
	} else {
		if start.x < end.x && start.y < end.y {
			l.drawCorridor(start, point{start.x, end.y})
			l.drawCorridor(point{start.x, end.y}, end)
		}
	}
}

func (l *level) draw() {
	var p point
	for p.x = 0; p.x < l.width; p.x++ {
		for p.y = 0; p.y < l.height; p.y++ {
			var c glyph
			tile := l.tile(p)
			if tile != nil /* && l.hasVisited(p) */ {
				c = tile.glyph()
			} else {
				c = ' '
			}
			drawGlyph(p, c, termbox.ColorDefault, termbox.ColorBlack)
		}
	}
}

func (l *level) drawRooms() []*room {
	tree := buildRooms(&room{0, 0, l.width - 1, l.height - 1}, 4)
	candidates := tree.rooms()

	var rooms []*room
	for _, c := range candidates {
		r := resizeRoom(c)
		if r == nil {
			continue
		}
		var ratio float64
		if r.w <= r.h {
			ratio = float64(r.w) / float64(r.h)
		} else {
			ratio = float64(r.h) / float64(r.w)
		}
		log.Printf("ratio %f %v\n", ratio, ratio < .5)
		if ratio > .25 {
			rooms = append(rooms, r)
		}
	}

	for _, r := range rooms {
		r.draw(l)
	}
	return rooms
}

func (l *level) drawCorridors(rooms []*room) {
	l.drawCorridor(point{3, 5}, point{20, 5})
	l.drawCorridor(point{5, 3}, point{5, 20})
	l.drawCorridor(point{20, 10}, point{3, 10})
	l.drawCorridor(point{10, 20}, point{10, 3})

	l.drawCorridor(point{30, 5}, point{50, 25})
}

func makeLevel(w, h int) *level {
	var l = level{
		width:      w,
		height:     h,
		layout:     make([]tiler, w*h),
		attributes: make([]attr, w*h),
	}

	l.drawCorridors(nil)
	log.Println("Done with make")

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
