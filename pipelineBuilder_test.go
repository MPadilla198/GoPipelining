package PipinHot

import "testing"

func assertAddStagePanic(t *testing.T, fn func(Function, uint) *builder, in Function, n uint, fnNum int) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AddStage method %d did not panic.", fnNum)
		}
	}()
	fn(in, n)
}

var b PipelineBuilder

func init() {
	b = NewPipelineBuilder()
}

func TestBuilder_AddStage(t *testing.T) {
	// TODO create panic-needed test cases
	// assertAddStagePanic(t, b.AddStage)
}

func TestBuilder_Build(t *testing.T) {

}

func TestBuilder(t *testing.T) {

}

func BenchmarkBuilder_AddStage(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkBuilder_Build(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkNewPipelineBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
