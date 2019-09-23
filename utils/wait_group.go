package utils

import (
	"sync"
	"sync/atomic"
)

type WaitGroup interface {
	Add(int)
	Done()
	Count() int
	Wait()
}

func newWaitGroup() WaitGroup {
	return &waitGroup{sync.WaitGroup{}, 0}
}

type waitGroup struct {
	wg  sync.WaitGroup
	cnt int64
}

func (wg *waitGroup) Add(delta int) {
	wg.wg.Add(delta)
	atomic.AddInt64(&wg.cnt, int64(delta))
}

func (wg *waitGroup) Done() {
	wg.wg.Done()
	atomic.AddInt64(&wg.cnt, -1)
}

func (wg *waitGroup) Count() int {
	return int(wg.cnt)
}

func (wg *waitGroup) Wait() {
	wg.wg.Wait()
}
