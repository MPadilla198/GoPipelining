package PipinHot

import (
	"errors"
	"github.com/MPadilla198/PipinHot/utils"
	"reflect"
	"sync"
)

type Pipeline interface {
	start()

	Execute(...interface{}) error
	Next() (interface{}, bool)
	Flush() []interface{}
	// pipeline panics if pipeline is used again after calling Close()
	Close()
}

func buildPipeline(stages []builderStage) Pipeline {
	inType := stages[0].inputType
	inChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, inType), 0)
	outChanType := stages[len(stages)-1].outputType
	newPipeline := &pipeline{inType, inChan, reflect.Zero(outChanType), reflect.MakeChan(reflect.ChanOf(reflect.BothDir, done), 0), sync.WaitGroup{}, make([]stageDispatcher, len(stages)), 0, make([]interface{}, 0)}

	for i, stage := range stages {
		newPipeline.stageDispatchers[i] = newStageDispatcher(stage)
		inChan = newPipeline.stageDispatchers[i].Start(inChan)
	}

	// inChan at this point will be the final outChan from the last piece of the pipeline
	newPipeline.outChan = inChan

	newPipeline.start()

	return newPipeline
}

type pipeline struct {
	inputType reflect.Type

	inChan  reflect.Value
	outChan reflect.Value

	endPipeline reflect.Value

	wg sync.WaitGroup

	stageDispatchers []stageDispatcher
	itemsInPipeline  utils.Counter

	// TODO Change to different data structure that is more suitable
	values []interface{}
}

func (p *pipeline) start() {
	go func() {
		for {
			chosen, recv, _ := reflect.Select([]reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: p.outChan},
				{Dir: reflect.SelectRecv, Chan: p.endPipeline},
			})
			switch chosen {
			case 0:
				p.values = append(p.values, []interface{}{recv.Interface()})
			case 1:
				return
			}
		}
	}()
}

func (p *pipeline) Execute(vals ...interface{}) error {
	// First needs to check that vals are the same type as input
	for _, val := range vals {
		if reflect.TypeOf(val) != p.inputType {
			return errors.New("")
		}
	}

	for _, input := range vals {
		p.inChan.Send(reflect.ValueOf(input))
		p.itemsInPipeline.Increment()
	}

	return nil
}

// TODO IMPLEMENT
func (p *pipeline) Next() (interface{}, bool) {
	return nil, false
}

// TODO IMPLEMENT
func (p *pipeline) Flush() []interface{} {
	return []interface{}{}
}

func (p *pipeline) Close() {
	for _, s := range p.stageDispatchers {
		s.Close()
	}

	// s.Close() closes all input channels
	// closes last channel left
	p.inChan.Close()

	p.endPipeline.Close()
}
