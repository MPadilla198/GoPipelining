package utils

import "sync/atomic"

/*
Going to be used to count goroutines to keep track of them
in automatic stage dispatchers

Don't think there'll be more goroutines than 2^32
*/
type Counter int32

func (c *Counter) Increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *Counter) Decrement() int32 {
	return atomic.AddInt32((*int32)(c), -1)
}

func (c *Counter) Get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func (c *Counter) Zero() {
	atomic.StoreInt32((*int32)(c), 0)
}
