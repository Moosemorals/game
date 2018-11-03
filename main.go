package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

type logger struct {
	lines   []string
	display int
}

func (l *logger) log(msg string) {
	l.lines = append(l.lines, msg)
}

func (l *logger) upTo() int {
	if l.display < len(l.lines) {
		return l.display
	}
	return len(l.lines)
}

func (l *logger) draw(x, y int) {
	for i := 0; i < l.upTo(); i++ {
		drawString(x, y+i, l.lines[len(l.lines)-(i+1)])
	}
}

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

type sprite struct {
	x, y int
	c    rune
}

func cap(min, max, value int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
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

func drawString(x, y int, msg string) {
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorBlue, termbox.ColorBlack)
	}
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

func main() {
	err := termbox.Init()
	if err != nil {
		log.Panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)

	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	var player = sprite{
		x: 5,
		y: 5,
		c: '@',
	}

	var l = makeLevel()

	l.draw()
	termbox.Flush()

	for e := range events {
		if e.Ch == 'q' {
			break
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		if e.Type == termbox.EventKey {
			player.handleKeyEvent(e, &l)
		}
		l.draw()
		player.draw()
		termbox.Flush()
	}
}
