package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/nsf/termbox-go"
)

type keyHandler interface {
	handleKeyEvent(termbox.Event, *context)
}

type context struct {
	size    point
	player  *sprite
	level   *level
	loggger *logger
}

func (c *context) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c.level.draw()
	c.player.draw()

	c.loggger.draw(point{1, c.size.y - (1 + c.loggger.display)})
	termbox.Flush()
}

func (c *context) handleKeyEvent(e termbox.Event) {
	c.level.handleKeyEvent(e, c)
	c.player.handleKeyEvent(e, c)
}

func (c *context) log(msg string) {
	c.loggger.log(msg)
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

func drawGlyph(at point, g glyph, fg, bg termbox.Attribute) {
	termbox.SetCell(at.x, at.y, rune(g), fg, bg)
}

func main() {
	rand.Seed(4)

	var size point
	//size.x, size.y = termbox.Size()
	size.x = 100
	size.y = 50

	context := context{
		size: size,
		player: &sprite{
			point: point{
				x: 5,
				y: 5,
			},
			c: '@',
		},
		level:   makeLevel(size.x, size.y),
		loggger: &logger{display: 3},
	}

	context.level.visit(context.player.point)

	context.log("Welcome!")
	if false {
		os.Exit(0)
	}

	log.Println("Initilizing termbox")
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
