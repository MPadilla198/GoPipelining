package PipinHot

type Pipeline interface {
	Execute(...interface{}) error
	Next() interface{}
	Flush() []interface{}
	Close()
}
