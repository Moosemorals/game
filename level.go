package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"math"
	"math/rand"
	"sort"
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

func (l *level) drawRoom(left, top, right, bottom int) {
	var p point
	for p.x = left; p.x <= right; p.x++ {
		for p.y = top; p.y <= bottom; p.y++ {
			if p.x == left || p.x == right || p.y == top || p.y == bottom {
				if p.x == left {
					if p.y == top {
						l.setTile(p, &wall{g: wallTopLeft})
					} else if p.y == bottom {
						l.setTile(p, &wall{g: wallBottomLeft})
					} else {
						l.setTile(p, &wall{g: wallVertical})
					}
				} else if p.x == right {
					if p.y == top {
						l.setTile(p, &wall{g: wallTopRight})
					} else if p.y == bottom {
						l.setTile(p, &wall{g: wallBottomRight})
					} else {
						l.setTile(p, &wall{g: wallVertical})
					}
				} else {
					l.setTile(p, &wall{g: wallHorizontal})
				}
			} else {
				l.setTile(p, new(floor))
			}
		}
	}
	/*
		var dpoint = point{right, (bottom - top) / 2}
		l.setTile(dpoint, &door{
			open:       false,
			horizontal: false,
			point:      dpoint,
		})
	*/
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

// Generation code loosley based on https://eskerda.com/bsp-dungeon-generation/
type box struct {
	x, y, w, h int
}

func (b *box) tl() point {
	return point{b.x, b.y}
}
func (b *box) center() point {
	return point{b.x + b.w/2, b.y + b.h/2}
}

type node struct {
	box         *box
	left, right *node
}

func (n *node) boxes() []*box {
	if n.left != nil && n.right != nil {
		return append(n.left.boxes(), n.right.boxes()...)
	}
	return []*box{n.box}
}

func random(min, max int) int {
	if (max - min) == 0 {
		return 0
	}
	return rand.Intn(max-min) + min
}

func splitBox(b *box, vh bool) (l, r *box) {
	var split int
	if vh {
		factor := float64(b.w) / (rand.Float64() + 1.5)
		split = int(math.Round(factor))

		// Vertical split
		l = &box{
			x: b.x,
			y: b.y,
			w: split,
			h: b.h,
		}
		r = &box{
			x: b.x + l.w,
			y: b.y,
			w: b.w - l.w,
			h: b.h,
		}
	} else {
		factor := float64(b.h) / (rand.Float64() + 1.5)
		split = int(math.Round(factor))

		// Horizontal split
		l = &box{
			x: b.x,
			y: b.y,
			w: b.w,
			h: split,
		}
		r = &box{
			x: b.x,
			y: b.y + l.h,
			w: b.w,
			h: b.h - l.h,
		}
	}
	return
}

func buildTree(b *box, depth int) *node {
	root := &node{box: b}
	if depth > 0 {
		leftBox, rightBox := splitBox(b, depth%2 == 0)
		root.left = buildTree(leftBox, depth-1)
		root.right = buildTree(rightBox, depth-1)
	}
	return root
}

func boxToRoom(b *box) *box {
	var r box
	tries := 3
	for {
		r.x = b.x + random(0, b.x/3)
		r.y = b.y + random(0, b.y/3)
		r.w = b.w - (r.x - b.x)
		r.h = b.h - (r.y - b.y)
		if r.w > 3 && r.h > 3 {
			return &r
		}
		tries--
		if tries == 0 {
			return nil
		}
	}
}

func orientation(p, q, r point) int {
	o := (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
	if o == 0 {
		return 0
	} else if 0 > 0 {
		return 1
	}
	return 2
}

func calcHull(rooms []*box) []*box {

	// Find room with smallest y
	index := 0
	minY := rooms[0].center().y
	for i := 1; i < len(rooms); i++ {
		r := rooms[i]
		if r.center().y < minY || (r.center().y == minY && r.center().x < rooms[index].center().x) {
			minY = r.center().y
			index = i
		}
	}

	// Swap rooms[0] with point with the lowest y
	if index > 0 {
		rooms[0], rooms[index] = rooms[index], rooms[0]
	}

	p := rooms[0]
	p0 := p.center()
	hull := rooms[1:]

	// Sort points by polar angle with points[0]
	sort.SliceStable(hull, func(i, j int) bool {
		left := rooms[i].center()
		right := rooms[j].center()
		o := orientation(p0, left, right)
		if o == 0 {
			return p0.dist(right) >= p0.dist(left)
		}
		return o == 2
	})

	hull = append([]*box{p}, hull...)

	m := 1
	for i := 1; i < len(hull); i++ {
		for i < len(hull)-1 && orientation(p0, hull[i].center(), hull[i+1].center()) == 0 {
			i++
		}
		hull[m] = hull[i]
		m++
	}

	if m < 3 {
		log.Panic("Not enough rooms")
	}

	var stack boxStack
	stack.push(hull[0])
	stack.push(hull[1])
	stack.push(hull[2])

	for i := 3; i < m; i++ {
		for orientation(stack.peek(1).center(), stack.peek(0).center(), hull[i].center()) != 2 {
			stack.pop()
		}
		stack.push(hull[i])
	}

	log.Println("stack ", stack)
	log.Println("rooms ", rooms)
	return stack
}

func makeLevel(w, h int) *level {
	var l = level{
		width:      w,
		height:     h,
		layout:     make([]tiler, w*h),
		attributes: make([]attr, w*h),
	}

	tree := buildTree(&box{0, 0, w - 1, h - 1}, 5)
	boxes := tree.boxes()

	var rooms []*box
	for _, b := range boxes {
		r := boxToRoom(b)
		if r == nil {
			continue
		}
		ratio := float64(r.w) / float64(r.h)
		log.Printf("ratio %f\n", ratio)
		if ratio < .7 || ratio > 4 {
			continue
		}
		rooms = append(rooms, r)
	}

	//	hull := calcHull(rooms)

	for _, r := range rooms {
		l.drawRoom(r.x, r.y, r.x+r.w-1, r.y+r.h-1)
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
