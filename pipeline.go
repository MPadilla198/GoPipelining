package PipinHot

import "errors"

type Pipeline interface {
	Execute(...interface{}) error
	Next() (interface{}, bool)
	WaitAndFlush() []interface{}
	// pipeline panics if pipeline is used again after calling Close()
	Close()
}

type pipeline struct {
	itemsInPipeline int
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
func (p *pipeline) WaitAndFlush() []interface{} {
	return []interface{}{}
}

// TODO IMPLEMENT
func (p *pipeline) Close() {
}

func newPipeline() Pipeline {
	return &pipeline{}
}
