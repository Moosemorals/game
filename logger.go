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

func (l *logger) draw(p point) {
	for i := 0; i < l.upTo(); i++ {
		drawString(p.add(point{0, i}), l.lines[len(l.lines)-(i+1)])
	}
}
