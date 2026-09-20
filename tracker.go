package tracer

import (
	"sync"

	"github.com/nsonidotdev/gotrail/internal/id"
)

type tracker struct {
	// map of span pointers by id key
	activeSpans map[id.SpanID]*Span
	mu          sync.Mutex
}

func newTracker() *tracker {
	return &tracker{
		activeSpans: make(map[id.SpanID]*Span),
	}
}

func (t *tracker) recordStart(s *Span) {
	t.mu.Lock()
	defer t.mu.Unlock()
	existing := t.activeSpans[s.ID]
	if existing != nil {
		return
	}

	t.activeSpans[s.ID] = s
}

func (t *tracker) recordFinish(id id.SpanID) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.activeSpans, id)
}

func (t *tracker) terminateAll(reason string) {
	t.mu.Lock()
	activeIDs := make([]id.SpanID, 0, len(t.activeSpans))
	for id := range t.activeSpans {
		activeIDs = append(activeIDs, id)
	}
	t.mu.Unlock()

	finishOptions := finishOptions{status: statusFail, reason: reason}
	for _, ID := range activeIDs {
		finishByID(ID, finishOptions)
	}
}
