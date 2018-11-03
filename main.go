package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

type keyHandler interface {
	handleKeyEvent(termbox.Event, *context)
}

type context struct {
	player *sprite
	level  *level
}

func (c *context) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c.level.draw()
	c.player.draw()
	termbox.Flush()
}

func (c *context) handleKeyEvent(e termbox.Event) {
	c.level.handleKeyEvent(e, c)
	c.player.handleKeyEvent(e, c)
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

	var context context

	context.player = &sprite{
		point: point{
			x: 5,
			y: 5,
		},
		c: '@',
	}

	context.level = makeLevel()

	context.draw()

	for e := range events {
		if e.Ch == 'q' {
			break
		}
		if e.Type == termbox.EventKey {
			context.handleKeyEvent(e)
		}
		context.draw()
	}
}
