package utils

import (
	"math"
	"sync"
	"time"
)

type ComputeDuration func(totalTime int64, times []int64) time.Duration
type StopTimer func() time.Duration

func Av() ComputeDuration {
	return func(totalTime int64, times []int64) time.Duration {
		return time.Duration(totalTime / int64(len(times)))
	}
}

func Std(n float64) ComputeDuration {
	return func(totalTime int64, times []int64) time.Duration {
		mean := float64(totalTime) / float64(len(times))

		var variance float64

		for _, t := range times {
			variance += math.Pow(float64(t)-mean, 2.0)
		}

		variance /= float64(len(times))

		// calculates std
		return time.Duration(mean + (math.Sqrt(variance) - n))
	}
}

type Timer interface {
	// Starts timer and returns StopTimer function that when called ends timer and returns computed duration
	Start() StopTimer

	// Returns previously computed duration
	Get() time.Duration
}

func NewTimer(cap int, placeholderDuration time.Duration, duration ComputeDuration) Timer {
	return &timer{
		totalTime:   0,
		times:       make([]int64, 0, cap),
		fn:          duration,
		mux:         sync.Mutex{},
		savedReturn: placeholderDuration,
	}
}

type timer struct {
	totalTime int64
	times     []int64

	fn  ComputeDuration
	mux sync.Mutex

	savedReturn time.Duration
}

func (t *timer) Start() StopTimer {
	var startTime time.Time

	endFunc := func() time.Duration {
		timeLapsed := int64(time.Since(startTime))

		t.mux.Lock()
		defer t.mux.Unlock()

		t.totalTime += timeLapsed
		t.times = append(t.times, timeLapsed)

		if cap(t.times) == len(t.times) {
			t.savedReturn = t.fn(t.totalTime, t.times)
			t.times = t.times[:0]
			t.totalTime = 0
		}

		return t.savedReturn
	}

	startTime = time.Now()

	return endFunc
}

func (t *timer) Get() time.Duration {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.savedReturn
}
