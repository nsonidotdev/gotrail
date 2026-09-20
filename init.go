package tracer

import (
	"sync"
	"sync/atomic"
)

type tracerConfig struct {
	service string
	tracker *tracker
}

var (
	initOnce      sync.Once
	tracer        *tracerConfig
	isInitialized atomic.Bool
	isTerminated  atomic.Bool
)

func Init(serviceName string) {
	initOnce.Do(func() {
		tracker := newTracker()
		tracer = &tracerConfig{tracker: tracker, service: serviceName}
		isInitialized.Store(true)
	})
}

// Shutdown terminates all traces. Call this for graceful shutdown
func Shutdown() {
	isTerminated.Store(true)

	if tracer != nil && tracer.tracker != nil {
		tracer.tracker.terminateAll(errProcTerminated.Error())
	}
}
