package main

import (
	"log"
)

type roomStack []*room

func (s *roomStack) push(b *room) {
	*s = append(*s, b)
}

func (s *roomStack) pop() *room {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *roomStack) peek(n int) *room {
	index := len(*s) - (n + 1)
	log.Printf("peek len %d %d %d\n", len(*s), n, index)
	return (*s)[index]
}
