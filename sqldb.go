package main

import (
	"database/sql"
	"fmt"
	"time"
)

//Provide a limit after which the circuit breaker trips
const breakLimit = 5

type sqlConnection interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type sqlDB struct{}

//Fake Query
func (db *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	time.Sleep(time.Millisecond * 100)
	return nil, nil
}

//Simulate concurrent DB requests
func makeConcurrentRequest(cb *circuitBreaker, ch chan<- string) {
	start := time.Now()
	_, err := cb.Execute()
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f elapsed with err: %s", secs, err)
}

func main() {
	cb := newCircuitBreaker(breakLimit)
	concurrent := 10
	ch := make(chan string)
	for i := 0; i < concurrent; i++ {
		go makeConcurrentRequest(cb, ch)
	}
	for i := 0; i < concurrent; i++ {
		fmt.Println(<-ch)
	}
}
