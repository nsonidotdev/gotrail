package gotrail

import (
	"sync"
	"time"

	"github.com/nsonidotdev/gotrail/internal/id"
)

type status string

const (
	statusRunning status = "running"
	statusSuccess status = "success"
	statusFail    status = "fail"
	statusSkip    status = "skip"
)

type Span struct {
	ID         id.SpanID
	Name       string
	Trace      *Trace
	Duration   time.Duration
	Start      time.Time
	Attributes map[string]any
	status     status
	reason     string
	children   []*Span
	parent     *Span
	mu         sync.Mutex
	// After func stop() callback
	stopCancelListener func() bool
}

func getRoot(s *Span) *Span {
	if s.parent == nil {
		return s
	}

	return getRoot(s.parent)
}
