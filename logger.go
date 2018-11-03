package main

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
