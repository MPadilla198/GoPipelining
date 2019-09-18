package PipinHot

import (
	"fmt"
	"github.com/MPadilla198/PipinHot/utils"
	"reflect"
	"time"
)

type stageDispatcher interface {
	Start(inChan reflect.Value) (outChan reflect.Value)
	newWorker()
	Close()
}

func newStageDispatcher(stage builderStage) stageDispatcher {
	// TODO Find optimal buffer size for out chan
	doneChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, done), 0)
	inChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.inputType), 0)
	outChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.outputType), 0)

	if stage.nodeCnt == 0 {
		intoFnChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.inputType), 0)

		return &automaticStageDispatcher{inChan, outChan, intoFnChan, stage.fn, doneChan, 0, 0}
	}
	return &manualStageDispatcher{inChan, outChan, stage.fn, doneChan, stage.nodeCnt}
}

type manualStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	// Function
	fn reflect.Value

	doneChan reflect.Value
	nodeCnt  uint
}

func (man *manualStageDispatcher) newWorker() {
	defer func() {
		if r, ok := recover().(string); ok {
			// Catch just this panic, otherwise keep going
			// TODO figure out why this panic is being set of when pipeline.Next() is called
			if r == "send on closed channel" {
				fmt.Println("recovered: send on closed channel")
				return
			}

			panic(r)
		}
	}()

	for {
		// Select from input of channels: in and done
		chosen, recv, _ := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: man.inChan},
			{Dir: reflect.SelectRecv, Chan: man.doneChan},
		})
		switch chosen {
		case 0: // Something comes in the channel
			// Call fptr with input from in-channel as param
			// And send it through the output channel
			man.outChan.Send(man.fn.Call([]reflect.Value{recv})[0])
		case 1: // Done channel
			return
		}
	}
}

func (man *manualStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	// TODO Find optimal buffer size for out chan
	man.inChan = inChan

	for i := uint(0); i < man.nodeCnt; i++ {
		go man.newWorker()
	}

	return man.outChan
}

func (man *manualStageDispatcher) Close() {
	man.doneChan.Close()
	man.outChan.Close()
}

type automaticStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	intoFnChan reflect.Value
	fn         reflect.Value

	doneChan    reflect.Value
	nodeCounter utils.Counter
	itemInStage utils.Counter
}

func (auto *automaticStageDispatcher) newWorker() {
	defer func() {
		if r, ok := recover().(string); ok {
			// Catch just this panic, otherwise keep going
			// TODO figure out why this panic is being set of when pipeline.Next() is called
			if r == "send on closed channel" {
				fmt.Println("recovered: send on closed channel")
				return
			}

			panic(r)
		}
	}()

	for {
		chosen, recv, _ := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: auto.intoFnChan},
			// TODO Change next line to be more dynamic timer (thinking of making it based on average times to complete task for this stage)
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(1 * time.Second))},
			{Dir: reflect.SelectRecv, Chan: auto.doneChan},
		})

		switch chosen {
		// New value comes in
		case 0:
			auto.outChan.Send(auto.fn.Call([]reflect.Value{recv})[0])
			auto.itemInStage.Decrement()
		// Timer goes off and worker shuts down, or done chan ends goroutine
		case 1, 2:
			auto.nodeCounter.Decrement()
			return
		}
	}
}

func (auto *automaticStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	auto.inChan = inChan

	// TODO IMPLEMENT
	// This goroutine will be in control of the amount of workers in the stage
	go func() {
		for {
			chosen, recv, _ := reflect.Select([]reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: auto.inChan},
				{Dir: reflect.SelectRecv, Chan: auto.doneChan},
			})
			switch chosen {
			case 0:
				auto.itemInStage.Increment()

				if auto.itemInStage.Get() > auto.nodeCounter.Get() {
					// Start Worker
					go auto.newWorker()

					auto.nodeCounter.Increment()
				}

				auto.intoFnChan.Send(recv)
			case 1:
				return
			}
		}
	}()

	return auto.outChan
}

func (auto *automaticStageDispatcher) Close() {
	auto.doneChan.Close()
	auto.outChan.Close()
}
