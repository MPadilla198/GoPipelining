package PipinHot

import (
	"github.com/MPadilla198/PipinHot/utils"
	"reflect"
)

type stageDispatcher interface {
	Start(inChan reflect.Value) (outChan reflect.Value)
	Close()
}

func newStageDispatcher(stage builderStage) stageDispatcher {
	// TODO Find optimal buffer size for out chan
	doneChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, done), 0)
	inChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.inputType), 0)
	outChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, stage.outputType), 0)

	if stage.nodeCnt == 0 {
		return &automaticStageDispatcher{inChan, outChan, stage.fn, doneChan, 0}
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

func (man *manualStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	// TODO Find optimal buffer size for out chan
	man.inChan = inChan

	for i := uint(0); i < man.nodeCnt; i++ {
		go man.fn.Call([]reflect.Value{man.doneChan, man.inChan, man.outChan})
	}

	return man.outChan
}

func (man *manualStageDispatcher) Close() {
	man.doneChan.Close()
}

type automaticStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	fn reflect.Value

	doneChan    reflect.Value
	nodeCounter utils.Counter
}

func (auto *automaticStageDispatcher) Start(inChan reflect.Value) reflect.Value {
	auto.inChan = inChan

	// TODO IMPLEMENT -
	go func() {

	}()

	return auto.outChan
}

func (auto *automaticStageDispatcher) Close() {
	auto.doneChan.Close()
}
