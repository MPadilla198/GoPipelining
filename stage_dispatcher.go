package PipinHot

import "reflect"

type stageDispatcher interface {
	Start()
	Close()
}

type manualStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	done chan struct{}
}

func (man *manualStageDispatcher) Start() {

}

func (man *manualStageDispatcher) Close() {
	man.inChan.Close()
	close(man.done)
}

type automaticStageDispatcher struct {
	inChan  reflect.Value
	outChan reflect.Value

	done chan struct{}

	currNodeCnt uint
}

func (auto *automaticStageDispatcher) Start() {

}

func (auto *automaticStageDispatcher) Close() {
	auto.inChan.Close()
	close(auto.done)
}
