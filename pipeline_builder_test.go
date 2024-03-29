package PipinHot

import (
	"fmt"
	"testing"
)

func assertAddStagePanic(t *testing.T, fn func(uint, Function), in Function, n uint, fnNum int, mustPanic bool) {
	defer func() {
		if r := recover(); (r != nil) != mustPanic {
			t.Errorf("AddStage method #%d did not pass.", fnNum)
			t.Error(r)
		}
	}()
	fn(n, in)
}

var b2 PipelineBuilder

func init() {
	b2 = NewPipelineBuilder()
	b2.AddStage(0, func(n int) int { return n + 1 })
	b2.AddStage(10, func(i int) bool { return i > 50 })
}

func TestBuilder_AddStage(t *testing.T) {
	testCases := []struct {
		input     interface{}
		nodeCnt   uint
		mustPanic bool
	}{
		// MUST PANIC - input is not a function
		{35, 0, true},
		// MUST PANIC - input function doesn't have correct # of input params
		{func(n int, str string) bool { return false }, 0, true},
		// MUST PANIC - input function doesn't have correct # of input params
		{func() bool { return false }, 0, true},
		// MUST PANIC - input function doesn't have correct # of output params
		{func(n int) (bool, error) { return n == 0, nil }, 0, true},
		// MUST PANIC - input function doesn't have correct # of output params
		{func(n int) { fmt.Print(n) }, 0, true},

		// MUST NOT PANIC
		{func(n int) int { return n + 1 }, 0, false},
		{func(str string) bool { return str == "Hello, World!" }, 10, false},
	}

	for i, tCase := range testCases {
		b := NewPipelineBuilder()
		assertAddStagePanic(t, b.AddStage, tCase.input, tCase.nodeCnt, i, tCase.mustPanic)
	}
}

func TestBuilder_Build(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicked during build")
			t.Error(r)
		}
	}()

	pipe := b2.Build()

	if pipe == nil {
		t.Errorf("Error: Pipline is nil")
	}
}

func TestBuilder(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	nB := NewPipelineBuilder()

	nB.AddStage(0, func(n int) string {
		if n > 100 {
			return "Hello, World!"
		}

		return "Hello, Tiny Planet!"
	})
	nB.AddStage(20, func(str string) bool { return str == "Hello, World!" })
	nB.AddStage(15, func(b bool) (s struct {
		n   int
		str string
	}) {
		if b {
			s.n = 12
			s.str = "Bye."
		} else {
			s.n = 4325982
			s.str = "Bye bye."
		}

		return
	})
	nB.AddStage(9999, func(s struct {
		n   int
		str string
	}) int {
		return s.n
	})

	nB.Build()
}

func BenchmarkBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nB := NewPipelineBuilder()

		nB.AddStage(0, func(n int) string {
			if n > 100 {
				return "Hello, World!"
			}

			return "Hello, Tiny Planet!"
		})
		nB.AddStage(20, func(str string) bool { return str == "Hello, World!" })
		nB.AddStage(15, func(b bool) (s struct {
			n   int
			str string
		}) {
			if b {
				s.n = 12
				s.str = "Bye."
			} else {
				s.n = 4325982
				s.str = "Bye bye."
			}

			return
		})
		nB.AddStage(9999, func(s struct {
			n   int
			str string
		}) int {
			return s.n
		})

		nB.Build()
	}
}
