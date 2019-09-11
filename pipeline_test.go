package PipinHot

import "testing"

func newTestPipe(n uint) Pipeline {
	return NewPipelineBuilder().AddStage(func(n int) int { return n * n }, 2).AddStage(func(n int) int { return n * 2 }, n).Build()
}

func simpleIntPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n int) int { return n * n }, 0).Build()
}

func simpleStringPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n string) string { return n + "..." }, 0).Build()
}

func simpleBoolPipe() Pipeline {
	return NewPipelineBuilder().AddStage(func(n bool) bool { return !n }, 0).Build()
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
		{[]interface{}{"geg", "sgse", "egsaf"}, false},
		{[]interface{}{"fsef"}, false},
	}

	strPipe := simpleStringPipe()

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

}
