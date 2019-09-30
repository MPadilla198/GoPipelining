package utils

import (
	"sync"
	"testing"
	"time"
)

func getTestCases() []time.Duration {
	intCases := []int64{1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 100, 200, 300, 400, 500, 23, 53, 265, 875, 34, 65, 9, 35, 7, 24, 80, 453, 326, 999, 432, 64, 534} // in milliseconds

	durCases := make([]time.Duration, len(intCases))

	for i, c := range intCases {
		durCases[i] = time.Duration(c * int64(time.Millisecond))
	}

	return durCases
}

func TestTimer_Av(t *testing.T) {
	timer := NewTimer(10, 100*time.Millisecond, Av())

	testCases := getTestCases()
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(len(testCases))

	for i, c := range testCases {
		go func(caseNum int, cas time.Duration, tTimer Timer) {
			defer waitGroup.Done()

			done := tTimer.Start()
			time.Sleep(cas)
			av := done()

			t.Logf("Test case #%d: %v milliseconds", caseNum, av)
		}(i, c, timer)
	}

	waitGroup.Wait()
}

func TestTimer_Std(t *testing.T) {
	timer := NewTimer(10, 100*time.Millisecond, Std(0))

	testCases := getTestCases()
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(len(testCases))

	for i, c := range testCases {
		go func(caseNum int, cas time.Duration, tTimer Timer) {
			defer waitGroup.Done()

			done := tTimer.Start()
			time.Sleep(cas)
			std := done()

			t.Logf("Test case #%d: %v milliseconds", caseNum, std)
		}(i, c, timer)
	}

	waitGroup.Wait()
}

func TestTimer(t *testing.T) {

}

func BenchmarkTimer(b *testing.B) {
}
