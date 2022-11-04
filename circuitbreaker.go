package main

import (
	"database/sql"
	"fmt"
	"time"
)

func main() {
	cb := newCircuitBreaker(breakLimit)
	concurrent := 10
	ch := make(chan string)
	for i := 0; i < concurrent; i++ {
		// time.Sleep(time.Millisecond * 1)
		go makeConcurrentRequest(cb, ch)
	}
	for i := 0; i < concurrent; i++ {
		fmt.Println(<-ch)
	}
}

// Provide a limit after which the circuit breaker trips
const breakLimit = 5

type sqlConnection interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type sqlDB struct{}

// Fake Query
func (db *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	time.Sleep(time.Millisecond * 10)
	return nil, nil
}

// Simulate concurrent DB requests
func makeConcurrentRequest(cb *circuitBreaker, ch chan<- string) {
	start := time.Now()
	_, err := cb.Execute()
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f elapsed with err: %s", secs, err)
}

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

// Execute will execute the db operation as long as
// the CircuitBreaker is not over capacity
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
