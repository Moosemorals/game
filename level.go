package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
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

	tree := buildRooms(&room{0, 0, w - 1, h - 1}, 4)
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
		r.draw(&l)
	}
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
