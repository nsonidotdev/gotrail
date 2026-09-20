package tracer

import (
	"context"
	"time"

	"github.com/nsonidotdev/gotrail/internal/id"
)

func Start(parent context.Context, name string, attrs map[string]any) (context.Context, *Span) {
	if !isInitialized.Load() || isTerminated.Load() {
		return parent, nil
	}

	parentSpan, _ := getCtxSpan(parent)
	if parentSpan != nil && parentSpan.Trace.isFinished.Load() {
		return parent, nil
	}

	id, err := id.GenerateSpanID()
	if err != nil {
		return parent, nil
	}

	newSpan := &Span{
		ID:         id,
		Name:       name,
		Start:      time.Now(),
		Attributes: attrs,
		status:     statusRunning,
		parent:     parentSpan,
	}

	if parentSpan == nil {
		newSpan.Trace = newTrace(newSpan)
	} else {
		newSpan.Trace = parentSpan.Trace
	}

	tracer.tracker.recordStart(newSpan)

	if parentSpan != nil {
		appendSpanChild(parentSpan, newSpan)
	}

	ctx := setCtxSpan(parent, newSpan)
	stop := context.AfterFunc(ctx, func() {
		newSpan.finish(finishOptions{status: statusFail, reason: errSpanCtxCancelled.Error()})
	})
	newSpan.stopCancelListener = stop

	return ctx, newSpan
}

func appendSpanChild(s *Span, child *Span) {
	if s == nil || child == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.children = append(s.children, child)
}
