package gotrail

import (
	"sync/atomic"

	"github.com/nsonidotdev/gotrail/internal/id"
)

type Trace struct {
	ID         id.TraceID
	root       *Span
	isFinished atomic.Bool
}

func newTrace(root *Span) *Trace {
	// TODO: handle possible error
	id, _ := id.GenerateTraceID()

	return &Trace{
		ID:   id,
		root: root,
	}
}
