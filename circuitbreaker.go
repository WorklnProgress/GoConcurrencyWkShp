package main

import (
	"fmt"
)

type circuitBreaker struct {
	sqlDB
	queue     chan struct{}
	numCalled int
}

func newCircuitBreaker(limit int) *circuitBreaker {
	q := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		q <- struct{}{}
	}
	s := sqlDB{}
	return &circuitBreaker{sqlDB: s, queue: q, numCalled: 0}
}

func (cb *circuitBreaker) Execute() (rows interface{}, err error) {
	cb.numCalled++
	if !cb.canExecute() {
		return cb.overCapacity()
	}
	defer cb.doneExecuting()
	return cb.sqlDB.Query("SELECT FAKE FROM FAKE")
}

func (cb *circuitBreaker) overCapacity() (interface{}, error) {
	return nil, fmt.Errorf("no more connections in pool over capacity")
}

func (cb *circuitBreaker) canExecute() bool {
	select {
	case <-cb.queue:
		return true
	default:
	}
	return false
}

func (cb *circuitBreaker) doneExecuting() {
	cb.queue <- struct{}{}
}
