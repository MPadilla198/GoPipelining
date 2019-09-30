package PipinHot

import (
	"reflect"
)

/*

TODO - Rework builder to use cap instead of setNodeCnt, gonna make all pipelines automatic with caps instead on each stage

*/

// For a while it'll just use interface{}
type Function interface{}

var done = reflect.TypeOf(struct{}{})

type PipelineBuilder interface {
	Build() Pipeline
	AddStage(uint, Function)
}

type builderStage struct {
	fn      reflect.Value
	nodeCnt uint

	inputType  reflect.Type
	outputType reflect.Type
}

type builder struct {
	stages         []builderStage
	lastOutputType reflect.Type
}

func NewPipelineBuilder() PipelineBuilder {
	return &builder{stages: make([]builderStage, 0), lastOutputType: nil}
}

func (b *builder) Build() Pipeline {
	return buildPipeline(b.stages)
}

// AddStage expects fptr to be a pointer to a non-nil function
// setNodeCnt sets an exact amount of nodes to be instantiated
// If setNodeCnt is set to 0, the builderStage node cnt will be controlled automatically
func (b *builder) AddStage(setNodeCnt uint, fptr Function) {
	// fptr is a pointer to a function.
	fn := reflect.ValueOf(fptr)
	fnParams := fn.Type()

	// Makes sure input function has 1 arg and 1 return value only
	// Also checks that fptr is actually a function
	if fnParams.NumIn() != 1 || fnParams.NumOut() != 1 {
		panic("Invalid number of parameters/returns in function")
	}

	// Param types
	inType := fnParams.In(0)
	outType := fnParams.Out(0)

	if b.lastOutputType != nil {
		if b.lastOutputType != inType {
			panic("Stage's inputs don't match pipeline outputs")
		}
	}

	b.stages = append(b.stages, builderStage{fn, setNodeCnt, inType, outType})
	b.lastOutputType = outType
}
