package main

import (
	"log"
)

type boxStack []*box

func (s *boxStack) push(b *box) {
	*s = append(*s, b)
}

func (s *boxStack) pop() *box {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *boxStack) peek(n int) *box {
	index := len(*s) - (n + 1)
	log.Printf("peek len %d %d %d\n", len(*s), n, index)
	return (*s)[index]
}
