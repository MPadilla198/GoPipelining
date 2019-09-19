package utils

import (
	"math"
	"sync"
	"time"
)

type Timer interface {
	Start() func()
	Av() time.Duration
	Std(float64) time.Duration
}

func NewTimer(sampleSize int, placeholder time.Duration) Timer {
	return &timer{
		sampleSize: sampleSize,
		times:      newTimes(sampleSize),
		mux:        sync.RWMutex{},
		av:         int64(placeholder),
		std:        int64(placeholder),
	}
}

type timer struct {
	sampleSize int // used for cap size in *times, stays const.

	times *times

	mux sync.RWMutex
	av  int64
	std int64
}

func (t *timer) Start() (end func()) {
	var startTime time.Time

	end = func() {
		timeLapsed := time.Since(startTime)

		times, totalTime, recalculate := t.times.add(int64(timeLapsed))

		if recalculate {
			go func() {
				newAv := totalTime / int64(t.sampleSize)

				var variation float64
				for _, v := range times {
					variation += math.Pow(float64(v-newAv), 2)
				}

				variation /= float64(t.sampleSize)

				newStd := math.Sqrt(variation)

				t.mux.Lock()
				defer t.mux.Unlock()

				t.av = newAv
				t.std = int64(newStd)
			}()
		}
	}

	startTime = time.Now()

	return
}

func (t *timer) Av() time.Duration {
	t.mux.RLock()
	defer t.mux.RUnlock()

	return time.Duration(t.av)
}

func (t *timer) Std(n float64) time.Duration {
	t.mux.RLock()
	defer t.mux.RUnlock()

	return time.Duration(t.av + int64(float64(t.std)*n))
}

type times struct {
	mux       sync.Mutex
	totalTime int64
	times     []int64
}

func newTimes(cap int) *times {
	return &times{sync.Mutex{}, 0, make([]int64, 0, cap)}
}

// takes in new value to add
// Returns times array, sum of times in array, and if times reached cap and has been reset
func (t *times) add(n int64) ([]int64, int64, bool) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.times = append(t.times, n)
	t.totalTime += n

	if len(t.times) != cap(t.times) {
		return t.times, t.totalTime, false
	}

	// save and clear t.times[]
	times := t.times
	t.times = t.times[:] // Resets slice without touching underlying array

	// save and clear t.totalTime
	totalTime := t.totalTime
	t.totalTime = 0

	return times, totalTime, true
}
