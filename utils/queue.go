package utils

import (
	"sync"
)

type Queue interface {
	Queue(interface{})
	Pop() (interface{}, bool)
	Size() int
	List() []interface{}
	Clear()
}

func NewQueue() Queue {
	return newLinkedList()
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

func newLinkedList() *linkedList {
	nodePool := sync.Pool{
		New: func() interface{} {
			return &node{nil, nil}
		},
	}

	return &linkedList{nil, nil, 0, sync.Mutex{}, nodePool}
}

func (ll *linkedList) Queue(val interface{}) {
	newNode := ll.nodePool.Get().(*node)
	newNode.value = val

	ll.mux.Lock()
	defer ll.mux.Unlock()

	if ll.size == 0 {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.next = newNode
		ll.tail = ll.tail.next
	}

	ll.size++
}

func (ll *linkedList) Pop() (interface{}, bool) {
	ll.mux.Lock()
	defer ll.mux.Unlock()

	if ll.size == 0 {
		return nil, false
	}

	val := ll.head.value
	oldHead := ll.head

	ll.head = ll.head.next

	oldHead.next = nil
	oldHead.value = nil
	ll.nodePool.Put(oldHead)

	return val, true
}

func (ll *linkedList) Size() int {
	return ll.size
}

func (ll *linkedList) List() []interface{} {
	list := make([]interface{}, ll.size)

	ll.mux.Lock()
	defer ll.mux.Unlock()

	node := ll.head
	for i := 0; node != nil; i++ {
		list[i] = node
		node = node.next
	}

	return list
}

func (ll *linkedList) Clear() {
	ll.mux.Lock()
	defer ll.mux.Unlock()

	ll.head = nil
	ll.tail = nil
	ll.size = 0
}
