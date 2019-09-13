package PipinHot

import "reflect"

// For a while it'll just use interface{}
type Function interface{}

var done = reflect.TypeOf(struct{}{})

type PipelineBuilder interface {
	Build() Pipeline
	AddStage(uint, Function)
}

type stage struct {
	fn      reflect.Value
	isAuto  bool
	nodeCnt uint
}

type builder struct {
	stages []stage
}

func NewPipelineBuilder() PipelineBuilder {
	return &builder{stages: make([]stage, 0)}
}

func (b *builder) Build() Pipeline {
	return newPipeline()
}

// AddStage expects fptr to be a pointer to a non-nil function
// setNodeCnt sets an exact amount of nodes to be instantiated
// If setNodeCnt is set to 0, the stage node cnt will be controlled automatically
func (b *builder) AddStage(setNodeCnt uint, fptr Function) {
	// BIG TODO needs to check that new stage matches with stage before it (params and returns)

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

	// New Function Type made from function inputted
	newFuncType := reflect.FuncOf(
		[]reflect.Type{reflect.ChanOf(reflect.RecvDir, done), reflect.ChanOf(reflect.RecvDir, inType)},
		[]reflect.Type{reflect.ChanOf(reflect.RecvDir, outType)},
		false)

	newStageFn := reflect.MakeFunc(newFuncType, func(in []reflect.Value) []reflect.Value {
		doneChan := in[0]
		inChan := in[1]
		outChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, outType), 0)

		go func(newOut reflect.Value) {
			defer newOut.Close()

			for {
				// Select from input of channels: in and done
				chosen, recv, _ := reflect.Select([]reflect.SelectCase{
					{reflect.SelectRecv, inChan, reflect.ValueOf(0)},
					{reflect.SelectRecv, doneChan, reflect.ValueOf(0)},
				})
				switch chosen {
				case 0: // Something comes in the channel
					// Call fptr with input from in-channel as param
					// And send it through the output channel
					newOut.Send(fn.Call([]reflect.Value{recv})[0])
				case 1: // Done channel
					return
				}
			}
		}(outChan)

		return []reflect.Value{outChan}
	})

	b.stages = append(b.stages, stage{newStageFn, setNodeCnt == 0, setNodeCnt})
}
