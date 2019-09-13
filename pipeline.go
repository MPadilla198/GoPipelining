package PipinHot

import (
	"errors"
	"reflect"
)

type Pipeline interface {
	Execute(...interface{}) error
	Next() (interface{}, bool)
	Flush() []interface{}
	// pipeline panics if pipeline is used again after calling Close()
	Close()
}

func buildPipeline(stages []builderStage) Pipeline {
	inChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stages[0].inputType), 0)
	outChanType := stages[len(stages)-1].outputType
	newPipeline := &pipeline{inChan, reflect.Zero(outChanType), make([]stageDispatcher, len(stages)), 0}

	for i, stage := range stages {
		newPipeline.stageDispatchers[i] = newStageDispatcher(stage)
		inChan = newPipeline.stageDispatchers[i].Start(inChan)
	}

	// inChan at this point will be the final outChan from the last piece of the pipeline
	newPipeline.outChan = inChan

	return newPipeline
}

type pipeline struct {
	inChan  reflect.Value
	outChan reflect.Value

	stageDispatchers []stageDispatcher
	itemsInPipeline  uint
}

// TODO IMPLEMENT
func (p *pipeline) Execute(vals ...interface{}) error {
	return errors.New("not implemented")
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
}
