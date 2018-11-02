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

type sprite struct {
	x, y int
	c    rune
}

func (s *sprite) move(dx, dy int) {
	width, height := termbox.Size()
	s.x += dx
	if s.x < 0 {
		s.x = 0
	} else if s.x > width {
		s.x = width
	}
	s.y += dy
	if s.y < 0 {
		s.y = 0
	} else if s.y > height {
		s.y = height
	}
}

func (s *sprite) handleKeyEvent(e termbox.Event) {
	if e.Ch != 0 {
		return
	}
	switch e.Key {
	case termbox.KeyArrowLeft:
		s.move(-1, 0)
		return
	case termbox.KeyArrowRight:
		s.move(1, 0)
		return
	case termbox.KeyArrowUp:
		s.move(0, -1)
		return
	case termbox.KeyArrowDown:
		s.move(0, 1)
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

func drawBox(top, left, bottom, right int, c rune) {
	for x := left; x < right; x++ {
		termbox.SetCell(x, top, c, termbox.ColorBlue, termbox.ColorBlack)
		termbox.SetCell(x, bottom, c, termbox.ColorBlue, termbox.ColorBlack)
	}

	for y := top + 1; y < bottom; y++ {
		termbox.SetCell(left, y, c, termbox.ColorGreen, termbox.ColorBlack)
		termbox.SetCell(right-1, y, c, termbox.ColorGreen, termbox.ColorBlack)
	}

	termbox.Flush()
}

func drawMap() {
	drawString(5, 5, "Hello")
	drawString(12, 5, "World")
	drawBox(4, 4, 6, 18, '#')
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
		x: 9,
		y: 9,
		c: '@',
	}
	//var l = logger{display: 3}
	for e := range events {
		if e.Ch == 'q' {
			break
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		drawMap()
		if e.Type == termbox.EventKey {
			player.handleKeyEvent(e)
			player.draw()
		}
		termbox.Flush()
	}
}
