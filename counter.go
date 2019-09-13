package PipinHot

import "sync/atomic"

/*
Going to be used to count goroutines to keep track of them
in automatic stage dispatchers

Don't think there'll be more goroutines than 2^32
*/
type counter int32

func (c *counter) increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *counter) decrement() int32 {
	return atomic.AddInt32((*int32)(c), -1)
}

func (c *counter) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
