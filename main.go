package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

func cap(min, max, value int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}
func drawString(x, y int, msg string) {
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorBlue, termbox.ColorBlack)
	}
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
