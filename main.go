package main

import "github.com/nsf/termbox-go"
import "log"

func drawString(x, y int, msg string) {
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorBlue, termbox.ColorBlack)
	}
	termbox.Flush()
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

	drawString(5, 5, "Hello")
	drawString(12, 5, "World")

	drawBox(4, 4, 6, 18, '#')
	<-events
}
