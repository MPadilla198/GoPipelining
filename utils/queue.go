package utils

import (
	"sync"
)

type Queue interface {
	Push(interface{})
	Pop() interface{}
	Size() int
	List() []interface{}
}

func NewQueue() Queue {
	return nil
}

type node struct {
	next  *node
	value interface{}
}

type linkedList struct {
	head *node
	tail *node

	size int

	mux      sync.Mutex
	nodePool sync.Pool
}

func (ll *linkedList) Push(val interface{}) {

}

func (ll *linkedList) Pop() interface{} {
	return nil
}

func (ll *linkedList) Size() int {
	return 0
}

func (ll *linkedList) List() []interface{} {
	return []interface{}{}
}
