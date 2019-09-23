package PipinHot

import (
	"github.com/MPadilla198/PipinHot/utils"
	"reflect"
	"sync"
	"time"
)

type stageDispatcher interface {
	callFunc(value reflect.Value)
	Start(inChan reflect.Value) (outChan reflect.Value)
	startWorkers(n uint)
	Close()
}

func newStageDispatcher(stage builderStage) stageDispatcher {
	// TODO Find optimal buffer size for out chan, doing this for now
	doneChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, done), int(stage.nodeCnt))
	inChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.inputType), int(stage.nodeCnt))
	outChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.outputType), int(stage.nodeCnt))

	if stage.nodeCnt == 0 {
		intoFnChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.inputType), 0)

		return &automaticStageDispatcher{
			inChan:      inChan,
			outChan:     outChan,
			intoFnChan:  intoFnChan,
			fn:          stage.fn,
			timer:       utils.NewTimer(10, 1*time.Second),
			doneChan:    doneChan,
			nodeCounter: 0,
			itemInStage: 0,
		}
	}
	return &manualStageDispatcher{inChan, outChan, stage.fn, sync.WaitGroup{}, doneChan, stage.nodeCnt}
}

type manualStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	// Function
	fn reflect.Value

	wg sync.WaitGroup

	doneChan reflect.Value
	nodeCnt  uint
}

func (man *manualStageDispatcher) callFunc(recv reflect.Value) {
	toSend := man.fn.Call([]reflect.Value{recv})[0]
	defer man.wg.Done()

	// Meant to end the race condition of sending over a channel that could potentially be closed
	chosen, _, _ := reflect.Select([]reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: man.doneChan},
		{Dir: reflect.SelectDefault},
	})
	switch chosen {
	case 0:
		return
	case 1:
		man.outChan.Send(toSend)
	}
}

func (man *manualStageDispatcher) startWorkers(n uint) {
	for i := uint(0); i < n; i++ {
		go func() {
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
					man.wg.Add(1)
					man.callFunc(recv)
				case 1: // Done channel
					return
				}
			}
		}()
	}
}

func (man *manualStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	man.inChan = inChan

	man.startWorkers(man.nodeCnt)

	return man.outChan
}

func (man *manualStageDispatcher) Close() {
	man.doneChan.Close()

	man.wg.Wait()

	man.outChan.Close()
}

type automaticStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	intoFnChan reflect.Value
	fn         reflect.Value

	timer utils.Timer

	doneChan    reflect.Value
	nodeCounter utils.Counter
	itemInStage utils.Counter
}

func (auto *automaticStageDispatcher) callFunc(recv reflect.Value) {
	toSend := auto.fn.Call([]reflect.Value{recv})[0]

	chosen, _, _ := reflect.Select([]reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: auto.doneChan},
		{Dir: reflect.SelectDefault},
	})
	switch chosen {
	case 0:
		return
	case 1:
		auto.outChan.Send(toSend)
	}
}

func (auto *automaticStageDispatcher) startWorkers(n uint) {
	for i := uint(0); i < n; i++ {
		go func() {
			for {
				chosen, recv, _ := reflect.Select([]reflect.SelectCase{
					{Dir: reflect.SelectRecv, Chan: auto.intoFnChan},
					{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(auto.timer.Av()))},
					{Dir: reflect.SelectRecv, Chan: auto.doneChan},
				})

				switch chosen {
				// New value comes in
				case 0:
					endTimer := auto.timer.Start()
					auto.callFunc(recv)
					auto.itemInStage.Decrement()
					endTimer()
				// Timer goes off and worker shuts down, or done chan ends goroutine
				case 1, 2:
					auto.nodeCounter.Decrement()
					return
				}
			}
		}()

		auto.nodeCounter.Increment()
	}
}

func (auto *automaticStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	auto.inChan = inChan

	// This goroutine will be in control of the amount of workers in the stage
	go func() {
		selectCases := []reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: auto.inChan},
			{Dir: reflect.SelectRecv, Chan: auto.doneChan},
		}

		for {
			chosen, recv, _ := reflect.Select(selectCases)
			switch chosen {
			case 0:
				auto.itemInStage.Increment()

				if delta := auto.itemInStage.Get() - auto.nodeCounter.Get(); delta > 0 {
					// Start Worker
					auto.startWorkers(uint(delta))
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
