package main

import (
	"testing"
)

func TestCircuitBreakerCalled(t *testing.T) {
	cb := newCircuitBreaker(breakLimit)
	_, err := cb.Execute()
	if err != nil {
		t.Log(err)
	}
	if cb.numCalled != 1 {
		t.Fatalf("Error in calling circuitbreaker, expected %d, found %d", 1, cb.numCalled)
	}
}
