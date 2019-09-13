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
	p.outChan.Close()
}

func newPipeline(in, out reflect.Value) Pipeline {
	return &pipeline{in, out, make([]stageDispatcher, 0), 0}
}
