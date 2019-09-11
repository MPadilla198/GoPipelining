package PipinHot

import "testing"

func newTestPipe(n uint) Pipeline {
	return NewPipelineBuilder().AddStage(func(n int) int { return n * n }, 2).AddStage(func(n int) int { return n * 2 }, n).Build()
}

func simpleIntPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n int) int { return n * n }, 1).Build()
}

func simpleStringPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n string) string { return n + "..." }, 1).Build()
}

func simpleBoolPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n bool) bool { return !n }, 1).Build()
}

// Checks if array has the same values regardless of order
func softEqual(a []interface{}, b []interface{}) bool {
	newA := make([]interface{}, len(a))
	newB := make([]interface{}, len(b))

	copy(newA, a)
	copy(newB, b)

loop:
	for _, nA := range newA {
		for i, nB := range newB {
			if nA == nB {
				newB = append(newB[:i], newB[i+1:])
				continue loop
			}
		}

		return false
	}

	return true
}

func TestPipeline_Execute(t *testing.T) {
	intTestCases := []struct {
		input         []interface{}
		errorExpected bool
	}{
		{[]interface{}{23, 43, "Hello", true}, true},
		{[]interface{}{"Hello", false}, true},
		{[]interface{}{}, false},
		{[]interface{}{23}, false},
		{[]interface{}{1, 2, 3, 4, 5, 6, 7}, false},
	}

	intPipe := simpleIntPipe()
	defer intPipe.Close()

	for _, tCase := range intTestCases {
		if err := intPipe.Execute(tCase.input); (err != nil) != tCase.errorExpected {
			t.Error("Execute is giving errors for ints.")
		}
	}

	strTestCases := []struct {
		input         []interface{}
		errorExpected bool
	}{
		{[]interface{}{"Hello", "hef", 23, false}, true},
		{[]interface{}{12, false, struct{ n int }{1}}, true},
		{[]interface{}{"name", "bat", "hat"}, false},
		{[]interface{}{"base"}, false},
	}

	strPipe := simpleStringPipe()
	defer strPipe.Close()

	for _, tCase := range strTestCases {
		if err := strPipe.Execute(tCase.input); (err != nil) != tCase.errorExpected {
			t.Error("Execute is giving errors for strings")
		}
	}

	boolTestCases := []struct {
		input         []interface{}
		errorExpected bool
	}{
		{[]interface{}{true, true, false}, false},
		{[]interface{}{false}, false},
		{[]interface{}{}, false},
		{[]interface{}{true, "Hello", 13}, true},
		{[]interface{}{"World.", 45, 98, 542, "Hello, "}, true},
	}

	boolPipe := simpleBoolPipe()
	defer boolPipe.Close()

	for _, tCase := range boolTestCases {
		if err := boolPipe.Execute(tCase.input); (err != nil) != tCase.errorExpected {
			t.Error("Execute is giving errors for bools.")
		}
	}
}

func TestPipeline_Next(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	testCases := []struct {
		input         []interface{}
		expectedValue []int
	}{
		{[]interface{}{10}, []int{200}},
		{[]interface{}{1, 2, 3, 4, 5, 6, 7}, []int{2, 8, 18, 32, 50, 72, 98}},
	}

	pipe := newTestPipe(1)
	defer pipe.Close()

	for i, tCase := range testCases {
		err := pipe.Execute(tCase.input...)

		if err != nil {
			t.Error(err)
		}

		if val, ok := pipe.Next(); ok {
			if iVal, ok := val.(int); !ok || iVal != tCase.expectedValue[i] {
				t.Error("Error: Incorrect value received from pipeline.")
			}
		} else {
			t.Error("Error: No value received")
		}
	}
}

func TestPipeline_Flush(t *testing.T) {
	testCases := []struct {
		input          []interface{}
		expectedOutput []interface{}
	}{
		{[]interface{}{1}, []interface{}{2}},
		{[]interface{}{1, 2, 3, 4, 5, 6, 7}, []interface{}{2, 8, 18, 32, 50, 72, 98}},
		{[]interface{}{10, 20, 30}, []interface{}{200, 800, 1800}},
		{[]interface{}{5, 10, 15, 20}, []interface{}{50, 200, 450, 800}},
	}

	pipe := newTestPipe(1)
	defer pipe.Close()

	for _, tCase := range testCases {
		err := pipe.Execute(tCase.input)

		if err != nil {
			t.Error("Error in execute in Flush Test.")
		}

		results := pipe.WaitAndFlush()

		if equal := softEqual(tCase.expectedOutput, results); !equal {
			t.Error("Flush is messing up")
		}
	}
}

func TestPipeline_Close(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Close() didn't cause a panic.")
		}
	}()

	pipe := newTestPipe(1)

	pipe.Close()

	// WILL CAUSE PANIC
	err := pipe.Execute(10, 1)
	if err != nil {
		// Never supposed to get here
	}
}

func TestPipeline(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("TestPipeline panicked.")
		}
	}()

	pipe := newTestPipe(1)
	defer pipe.Close()

	results := make([]interface{}, 0)

	err := pipe.Execute(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	if err != nil {
		t.Error(err)
		return
	}

	v, ok := pipe.Next()
	if !ok {
		t.Error("Pipe Next() problems.")
	}
	results = append(results, v)

	flushed := pipe.WaitAndFlush()
	results = append(results, flushed)

	if isEqual := softEqual(results, []interface{}{2, 8, 18, 32, 50, 72, 98, 128, 162, 200}); !isEqual {
		t.Error("Results don't match!")
	}
}

func TestAutoPipeline(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("TestPipeline panicked.")
		}
	}()

	pipe := newTestPipe(1)
	defer pipe.Close()
	results := make([]interface{}, 0)

	err := pipe.Execute(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	if err != nil {
		t.Error(err)
		return
	}

	v, ok := pipe.Next()
	if !ok {
		t.Error("Pipe Next() problems.")
	}
	results = append(results, v)

	flushed := pipe.WaitAndFlush()
	results = append(results, flushed)

	if isEqual := softEqual(results, []interface{}{2, 8, 18, 32, 50, 72, 98, 128, 162, 200}); !isEqual {
		t.Error("Results don't match!")
	}
}

func BenchmarkPipeline(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Error("TestPipeline panicked.")
		}
	}()

	for i := 0; i < b.N; i++ {
		pipe := newTestPipe(1)
		results := make([]interface{}, 1)

		err := pipe.Execute(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
		if err != nil {
			b.Error(err)
			return
		}

		v, ok := pipe.Next()
		if !ok {
			b.Error("Pipe Next() problems.")
		}
		results = append(results, v)

		flushed := pipe.WaitAndFlush()
		results = append(results, flushed)

		if isEqual := softEqual(results, []interface{}{2, 8, 18, 32, 50, 72, 98, 128, 162, 200, 242, 288, 338, 392, 450, 512, 578, 648, 722, 800}); !isEqual {
			b.Error("Results don't match!")
		}

		pipe.Close()
	}
}

func BenchmarkAutoPipeline(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Error("TestPipeline panicked.")
		}
	}()

	for i := 0; i < b.N; i++ {
		pipe := newTestPipe(0)
		results := make([]interface{}, 0)

		err := pipe.Execute(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
		if err != nil {
			b.Error(err)
			return
		}

		v, ok := pipe.Next()
		if !ok {
			b.Error("Pipe Next() problems.")
		}
		results = append(results, v)

		flushed := pipe.WaitAndFlush()
		results = append(results, flushed)

		if isEqual := softEqual(results, []interface{}{2, 8, 18, 32, 50, 72, 98, 128, 162, 200, 242, 288, 338, 392, 450, 512, 578, 648, 722, 800}); !isEqual {
			b.Error("Results don't match!")
		}

		pipe.Close()
	}
}
