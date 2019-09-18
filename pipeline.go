package PipinHot

import (
	"errors"
	"github.com/MPadilla198/PipinHot/utils"
	"reflect"
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

	waitingForNextChan := make(chan interface{})

	newPipeline := &pipeline{inType, inChan, reflect.Zero(outChanType), reflect.MakeChan(reflect.ChanOf(reflect.BothDir, done), 0), make([]stageDispatcher, len(stages)), 0, false, waitingForNextChan, utils.NewQueue()}

	for i, stage := range stages {
		newPipeline.stageDispatchers[i] = newStageDispatcher(stage)
		inChan = newPipeline.stageDispatchers[i].Start(inChan)
	}

	// inChan at this point will be the final outChan from the last piece of the pipeline
	newPipeline.outChan = inChan

	newPipeline.start()

	return newPipeline
}

// THIS PIPELINE IS NOT THREAD SAFE, TO BE USED ONLY IN MAIN THREAD/ONE GOROUTINE
type pipeline struct {
	// This is the type for the input channel
	inputType reflect.Type

	// The channels for going into and coming out of the pipeline
	inChan  reflect.Value
	outChan reflect.Value

	// The done channel, close this when all components of pipeline are closed to finish closing last of pipeline
	endPipeline reflect.Value

	// TODO find possible way to combine functionality of wg and itemsInPipeline
	// wg sync.WaitGroup

	// Dispatchers control # of nodes in a stage, just need to be started to open up pipeline
	stageDispatchers []stageDispatcher
	// Counters the amount of items IN THE PIPELINE CURRENTLY
	// As soon as items come out of outChan decrement counter
	itemsInPipeline utils.Counter

	// Used for waiting for next item in pipeline
	waitingForNext bool
	// When waiting for next value, the value is sent through nextChan
	nextChan chan interface{}

	// Stores values coming out of pipeline
	values utils.Queue
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
				if p.waitingForNext {
					p.nextChan <- recv.Interface()
					p.waitingForNext = false
				} else {
					p.values.Queue(recv.Interface())
				}

				p.itemsInPipeline.Decrement()
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

	// Will panic if Execute() is called while Flush hasn't returned

	for _, input := range vals {
		p.inChan.Send(reflect.ValueOf(input))
		p.itemsInPipeline.Increment()
	}

	return nil
}

func (p *pipeline) Next() (interface{}, bool) {
	if p.values.Size() == 0 {
		if p.itemsInPipeline.Get() == 0 {
			return nil, false
		}

		p.waitingForNext = true
		select {
		case nextVal := <-p.nextChan:
			return nextVal, true
		}
	}

	return p.values.Pop()
}

func (p *pipeline) Flush() []interface{} {
	for p.itemsInPipeline > 0 {
	}

	list := p.values.List()
	p.values.Clear()

	return list
}

func (p *pipeline) Close() {
	// s.Close() closes all input channels
	// closes in channel so no new
	p.inChan.Close()

	// Close all dispatchers
	for _, s := range p.stageDispatchers {
		s.Close()
	}

	// Finally, close pipeline
	p.endPipeline.Close()
}
