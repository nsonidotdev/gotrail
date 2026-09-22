package gotrail

import (
	"sync"
	"sync/atomic"
)

type tracerConfig struct {
	service   string
	tracker   *tracker
	connector Connector
}

var (
	initOnce      sync.Once
	tracer        *tracerConfig
	isInitialized atomic.Bool
	isTerminated  atomic.Bool
)

func Init(service string, opts ...Option) {
	initOnce.Do(func() {
		tracker := newTracker()
		tracer = &tracerConfig{tracker: tracker, service: service}

		for _, opt := range opts {
			opt(tracer)
		}

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
