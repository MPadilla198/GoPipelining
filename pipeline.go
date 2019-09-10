package PipinHot

type Pipeline interface {
	Execute(...interface{}) error
	Next() interface{}
	Flush() []interface{}
	Close()
}

type pipeline struct{}

// TODO IMPLEMENT
func (p *pipeline) Execute(vals ...interface{}) error {
	return nil
}

// TODO IMPLEMENT
func (p *pipeline) Next() interface{} {
	return 0
}

// TODO IMPLEMENT
func (p *pipeline) Flush() []interface{} {
	return []interface{}{}
}

// TODO IMPLEMENT
func (p *pipeline) Close() {
}

func newPipeline() Pipeline {
	return &pipeline{}
}
